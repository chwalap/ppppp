package shared

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// WeatherPort is a port to weather service
var WeatherPort = 5001

// WeatherEndpoint point to weather endpoint
var WeatherEndpoint = fmt.Sprintf("http://localhost:%d/weather", WeatherPort)

// WeatherStatus holds timestamp and temperature
type WeatherStatus struct {
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

// WeatherResponse returns weather status
type WeatherResponse struct {
	Status WeatherStatus `json:"status"`
	Err    string        `json:"err,omitempty"`
}

// Database is a database handle
type Database struct {
	*sql.DB
}

// Db is global dabatase object
var Db *Database

func init() {
	db, err := sql.Open("sqlite3", "./weather.db")
	if err != nil {
		panic(err)
	}
	Db = &Database{db}
}

// CreateDbIfNotExists creates database
func (d *Database) CreateDbIfNotExists() error {
	query, err := d.Prepare(`
		CREATE TABLE IF NOT EXISTS Weather (
			id INTEGER PRIMARY KEY,
			timestamp DATETIME,
			temperature REAL
		)
	`)
	if err != nil {
		return err
	}
	_, err = query.Exec()
	return err
}

// WriteWeatherStatus writes weather status to db
func (d *Database) WriteWeatherStatus(status WeatherStatus) error {
	query, err := d.Prepare(fmt.Sprintf(`
	INSERT INTO Weather (timestamp, temperature) VALUES (
		datetime(%d, 'unixepoch'), %f
	)`, status.Timestamp, status.Temperature))
	if err != nil {
		return err
	}
	_, err = query.Exec()
	return err
}

// ReadWeatherStatuses reads weather status from db
func (d *Database) ReadWeatherStatuses() ([]WeatherStatus, error) {
	rows, err := d.Query(`
	SELECT temperature, CAST(strftime('%s', timestamp) AS INT) AS timestamp FROM Weather
	`)
	if err != nil {
		return []WeatherStatus{}, err
	}

	var timestamp int64
	var temperature float64
	statuses := []WeatherStatus{}
	for rows.Next() {
		rows.Scan(&temperature, &timestamp)
		statuses = append(statuses, WeatherStatus{Temperature: temperature, Timestamp: timestamp})
	}

	return statuses, nil
}
