package shared

import (
	"database/sql"
	"fmt"

	// import sqlite3 driver with init
	_ "github.com/mattn/go-sqlite3"
)

// Database is a database handle
type Database struct {
	*sql.DB
}

// Db is global dabatase object
var Db *Database

func init() {
	db, err := sql.Open("sqlite3", "/db/weather.db")
	if err != nil {
		panic(err)
	}
	Db = &Database{db}
}

// CreateDbIfNotExists creates database structure
func (d *Database) CreateDbIfNotExists(clean bool) error {
	var err error

	if clean {
		var name string
		tnames := []string{}
		if rows, err := d.Query(`SELECT name FROM sqlite_master WHERE type='table'`); err == nil {
			for rows.Next() {
				if err = rows.Scan(&name); err == nil {
					tnames = append(tnames, name)
				} else {
					return err
				}
			}
		} else {
			return err
		}

		for _, name = range tnames {
			if _, err = d.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, name)); err != nil {
				return err
			}
		}
	}

	if _, err = d.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id       INTEGER      PRIMARY KEY NOT NULL,
			username VARCHAR(255) NOT NULL UNIQUE,
			password BLOB         NOT NULL
		);

		CREATE TABLE IF NOT EXISTS cities (
			id   INTEGER      NOT NULL,
			name VARCHAR(255) NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS workers (
			id       INTEGER PRIMARY KEY NOT NULL,
			userid   INTEGER NOT NULL,
			cityid   INTEGER NOT NULL,
			interval INTEGER NOT NULL,
			running  INTEGER NOT NULL DEFAULT 0,

			FOREIGN KEY(cityid) REFERENCES cities(id),
			FOREIGN KEY(userid) REFERENCES users(id)
		);
		
		CREATE TABLE IF NOT EXISTS weather (
			id          INTEGER PRIMARY KEY NOT NULL,
			workerid    INTEGER NOT NULL,
			timestamp   INTEGER NOT NULL,
			temperature REAL    NOT NULL,

			FOREIGN KEY(workerid) REFERENCES workers(id)
		);

		CREATE TABLE IF NOT EXISTS sessions (
			token  TEXT PRIMARY KEY NOT NULL,
			data   BLOB NOT NULL,
			expiry REAL NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS city_idx ON cities(id);
		CREATE INDEX IF NOT EXISTS workers_idx ON workers(id);
		CREATE INDEX IF NOT EXISTS workers_users_idx ON workers(userid);
		CREATE INDEX IF NOT EXISTS weather_time_idx ON weather(timestamp);
		CREATE INDEX IF NOT EXISTS weather_worker_idx ON weather(workerid);
		CREATE INDEX IF NOT EXISTS users_idx ON users(username);
		CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);
	`); err != nil {
		return err
	}

	return nil
}

// User defines user in db
type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password []byte `json:"password" db:"password"`
}

// AddUser adds new user into db
func (d *Database) AddUser(user User) (int64, error) {
	var err error
	if stmt, err := d.Prepare(fmt.Sprintf(`INSERT INTO users (username, password) VALUES ('%s', ?)`, user.Username)); err == nil {
		if result, err := stmt.Exec(user.Password); err == nil {
			return result.LastInsertId()
		}
	}
	return -1, err
}

// GetUserByID gets user from db based on userid
func (d *Database) GetUserByID(userid int64) (User, error) {
	var user User
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, username, password FROM users WHERE id=%d`, userid)); err == nil {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&user.ID, &user.Username, &user.Password)
		}
	}
	return user, err
}

// GetUserByName gets user from db based on username
func (d *Database) GetUserByName(username string) (User, error) {
	var user User
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, username, password FROM users WHERE username = '%s'`, username)); err == nil {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&user.ID, &user.Username, &user.Password)
		}
	}
	return user, err
}

// DoesUserExists check whether user alrady exists
func (d *Database) DoesUserExists(username string) (bool, error) {
	var exists bool = false
	var result int = 0
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT EXISTS( SELECT 1 FROM users WHERE username='%s' )`, username)); err == nil {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&result)
			if result == 1 {
				exists = true
			}
		}
	}
	return exists, err
}

// Worker defines worker in db
type Worker struct {
	ID       int64         `json:"id" db:"id"`
	UserID   int64         `json:"userid" db:"userid"`
	CityID   int64         `json:"cityid" db:"cityid"`
	Interval int           `json:"interval" db:"interval"`
	Running  int           `json:"running" db:"running"`
	Stop     chan struct{} `json:"-"`
	User     User          `json:"-"`
	City     City          `json:"-"`
}

// AddWorker adds new worker into db
func (d *Database) AddWorker(worker Worker) (int64, error) {
	var err error
	if result, err := d.Exec(fmt.Sprintf(`INSERT INTO workers (userid, cityid, interval, running) VALUES (%d, %d, %d, %d)`, worker.UserID, worker.CityID, worker.Interval, worker.Running)); err == nil {
		return result.LastInsertId()
	}
	return -1, err
}

// DeleteWorker deletes worker from db
func (d *Database) DeleteWorker(workerid int64) error {
	_, err := d.Exec(fmt.Sprintf(`DELETE FROM workers WHERE id=%d`, workerid))
	return err
}

// UpdateWorker updates worker
func (d *Database) UpdateWorker(worker Worker) error {
	_, err := d.Exec(fmt.Sprintf(`UPDATE workers SET cityid=%d, interval=%d, running=%d WHERE id=%d`, worker.CityID, worker.Interval, worker.Running, worker.ID))
	return err
}

// DeleteWorkerData deletes weather entries for worker
func (d *Database) DeleteWorkerData(worker Worker) error {
	_, err := d.Exec(fmt.Sprintf(`DELETE FROM weather WHERE workerid=%d`, worker.ID))
	return err
}

// GetWorkerByID gets worker based on workerid
func (d *Database) GetWorkerByID(workerid int64) (Worker, error) {
	var worker Worker
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, userid, cityid, interval, running FROM workers WHERE id=%d`, workerid)); err == nil {
		defer rows.Close()
		if rows.Next() {
			if err = rows.Scan(&worker.ID, &worker.UserID, &worker.CityID, &worker.Interval, &worker.Running); err == nil {
				if worker.User, err = d.GetUserByID(worker.UserID); err != nil {
					return worker, err
				}
				if worker.City, err = d.GetCityByID(worker.CityID); err != nil {
					return worker, err
				}
				worker.Stop = make(chan struct{})
				return worker, err

			}
		}
	}
	return worker, err
}

// GetWorkerByUserCity gets worker based on userid and cityid
func (d *Database) GetWorkerByUserCity(user User, city City) (Worker, error) {
	var worker Worker
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, userid, cityid, interval, running FROM workers WHERE userid=%d AND cityid=%d`, user.ID, city.ID)); err == nil {
		defer rows.Close()
		if rows.Next() {
			if err = rows.Scan(&worker.ID, &worker.UserID, &worker.CityID, &worker.Interval, &worker.Running); err == nil {
				worker.User = user
				worker.City = city
				worker.Stop = make(chan struct{})
			}
		}
	}
	return worker, err
}

// GetUserWorkers gets workers assigned to userid
func (d *Database) GetUserWorkers(user User) (map[int64]Worker, error) {
	var workers map[int64]Worker = map[int64]Worker{}
	var w Worker
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, userid, cityid, interval, running FROM workers WHERE userid=%d`, user.ID)); err == nil {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&w.ID, &w.UserID, &w.CityID, &w.Interval, &w.Running); err == nil {
				if w.City, err = d.GetCityByID(w.CityID); err != nil {
					return nil, err
				}
				w.User = user
				w.Stop = make(chan struct{})
				workers[w.ID] = w
			} else {
				return nil, err
			}
		}
	}
	return workers, err
}

// GetWorkers reads all workers from db
func (d *Database) GetWorkers() (map[int64]Worker, error) {
	var w Worker
	var workers map[int64]Worker = map[int64]Worker{}
	var err error
	if rows, err := d.Query(`SELECT id, userid, cityid, interval, running FROM workers`); err == nil {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&w.ID, &w.UserID, &w.CityID, &w.Interval, &w.Running); err == nil {
				if w.User, err = d.GetUserByID(w.UserID); err != nil {
					return nil, err
				}
				if w.City, err = d.GetCityByID(w.CityID); err != nil {
					return nil, err
				}
				w.Stop = make(chan struct{})
				workers[w.ID] = w
			} else {
				return nil, err
			}
		}
	}
	return workers, err
}

// City defines city in db
type City struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// AddCity adds new city into db
func (d *Database) AddCity(city City) (int64, error) {
	var err error
	if result, err := d.Exec(fmt.Sprintf(`INSERT OR IGNORE INTO cities (id, name) VALUES (%d, '%s')`, city.ID, city.Name)); err == nil {
		return result.LastInsertId()
	}
	return -1, err
}

// GetCityByID gets city from db based on cityid
func (d *Database) GetCityByID(cityid int64) (City, error) {
	var city City
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT id, name FROM cities WHERE id=%d`, cityid)); err == nil {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&city.ID, &city.Name)
		}
	}
	return city, err
}

// GetUserCities returns all cities for which user has workes
func (d *Database) GetUserCities(userid int64) (map[int64]City, error) {
	var c City
	var cities map[int64]City = map[int64]City{}
	var err error
	if rows, err := d.Query(fmt.Sprintf(`
		SELECT
			c.id AS id, c.name AS name
		FROM
			workers AS w
		JOIN
			cities AS c
		ON
			w.cityid=c.id
		WHERE
			w.userid=%d`,
		userid)); err == nil {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&c.ID, &c.Name); err == nil {
				cities[c.ID] = c
			} else {
				return nil, err
			}
		}
	}
	return cities, err
}

// WeatherStatus defines weather status in db
type WeatherStatus struct {
	WorkerID    int64   `json:"workerid"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

// AddWeatherStatus adds new weather status into db
func (d *Database) AddWeatherStatus(status WeatherStatus) (int64, error) {
	var err error
	if result, err := d.Exec(fmt.Sprintf(`INSERT INTO weather (workerid, timestamp, temperature) VALUES (%d, %d, %f)`, status.WorkerID, status.Timestamp, status.Temperature)); err == nil {
		return result.LastInsertId()
	}
	return -1, err
}

// GetGatheringStartEnd returns first and last weather timestamp based on workerid
func (d *Database) GetGatheringStartEnd(worker Worker) (int64, int64, error) {
	var min, max int64
	var err error
	if rows, err := d.Query(fmt.Sprintf(`SELECT MIN(timestamp) AS min, MAX(timestamp) AS max FROM weather WHERE workerid=%d`, worker.ID)); err == nil {
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&min, &max)
		}
	}
	return min, max, err
}

// ReadWeatherStatuses reads weather statuses from db based on timestamp range and workerid
func (d *Database) ReadWeatherStatuses(worker Worker, start, end int64) ([]WeatherStatus, error) {
	var statuses []WeatherStatus = []WeatherStatus{}
	var s WeatherStatus
	var err error
	if rows, err := d.Query(fmt.Sprintf(`
		SELECT
			temperature, timestamp
		FROM
			weather
		WHERE
			workerid=%d AND timestamp BETWEEN %d AND %d`,
		worker.ID, start, end)); err == nil {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&s.Temperature, &s.Timestamp); err == nil {
				statuses = append(statuses, s)
			} else {
				return nil, err
			}
		}
	}
	return statuses, err
}
