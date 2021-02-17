package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"shared"
	"worker/weather"
)

// C is map of cities
var C map[int64]shared.City

// W is slice of workers
var W map[int64]shared.Worker

// M is worker's mutex
var M *sync.RWMutex

func init() {
	var err error
	if C, err = shared.LoadCities(); err != nil {
		panic(err)
	}

	M = &sync.RWMutex{}
}

func httperr(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func getWorkerFromRequest(r *http.Request) (shared.Worker, error) {
	var err error
	var body []byte
	var worker shared.Worker

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return worker, err
	}
	err = json.Unmarshal(body, &worker)
	return worker, err
}

func respondSuccess(w http.ResponseWriter) error {
	var err error
	var jsonStr []byte
	response := shared.WorkerResponse{
		Success: true,
		Err:     "",
	}
	if jsonStr, err = json.Marshal(response); err != nil {
		return err
	}
	w.Write(jsonStr)
	w.WriteHeader(http.StatusOK)
	return nil
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}

	var err error
	var worker shared.Worker

	if worker, err = getWorkerFromRequest(r); err != nil {
		httperr(w, fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}
	if worker.Interval < 1 || worker.Interval > 86400 {
		httperr(w, fmt.Errorf("interval must be in range <1; 86400>, but is %d", worker.Interval), http.StatusBadRequest)
		return
	}
	if _, err = shared.Db.AddCity(C[worker.CityID]); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if worker.ID, err = shared.Db.AddWorker(worker); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(worker.ID); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Worker #%d created: { city: %s; user: %s; interval: %ds }\n", worker.ID, worker.City.Name, worker.User.Username, worker.Interval)

	M.Lock()
	W[worker.ID] = worker
	M.Unlock()

	go Worker(worker)

	respondSuccess(w)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}

	var err error
	var worker shared.Worker

	if worker, err = getWorkerFromRequest(r); err != nil {
		httperr(w, fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}

	M.Lock()
	if W[worker.ID].Running == 1 {
		W[worker.ID].Stop <- struct{}{}
	}
	delete(W, worker.ID)
	M.Unlock()

	if err = shared.Db.DeleteWorkerData(worker); err != nil {
		M.Unlock()
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if err = shared.Db.DeleteWorker(worker.ID); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Worker #%d deleted: { city: %s; user: %s; interval: %ds }\n", worker.ID, worker.City.Name, worker.User.Username, worker.Interval)

	respondSuccess(w)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad request method"), http.StatusBadRequest)
		return
	}

	var err error
	var worker shared.Worker

	if worker, err = getWorkerFromRequest(r); err != nil {
		httperr(w, fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}
	if worker.Interval < 1 || worker.Interval > 86400 {
		httperr(w, fmt.Errorf("interval must be in range <1; 86400>, but is %d", worker.Interval), http.StatusBadRequest)
		return
	}
	if _, err = shared.Db.AddCity(C[worker.CityID]); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if err = shared.Db.UpdateWorker(worker); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(worker.ID); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Worker #%d edited: { city: %s; user: %s; interval: %ds }\n", worker.ID, worker.City.Name, worker.User.Username, worker.Interval)

	M.Lock()
	if worker.CityID != W[worker.ID].CityID {
		if err = shared.Db.DeleteWorkerData(W[worker.ID]); err != nil {
			M.Unlock()
			httperr(w, err, http.StatusInternalServerError)
			return
		}
	}
	if W[worker.ID].Running == 1 {
		W[worker.ID].Stop <- struct{}{}
	}
	W[worker.ID] = worker
	M.Unlock()

	go Worker(worker)

	respondSuccess(w)
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad request method"), http.StatusBadRequest)
		return
	}

	var err error
	var worker shared.Worker

	if worker, err = getWorkerFromRequest(r); err != nil {
		httperr(w, err, http.StatusBadRequest)
		return
	}
	if W[worker.ID].Running == 0 {
		httperr(w, fmt.Errorf("worker is not running"), http.StatusBadRequest)
		return
	}
	worker.Running = 0
	if err = shared.Db.UpdateWorker(worker); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(worker.ID); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	M.Lock()
	W[worker.ID].Stop <- struct{}{}
	W[worker.ID] = worker
	M.Unlock()

	log.Printf("Worker #%d paused: { city: %s; user: %s; interval: %ds }\n", worker.ID, worker.City.Name, worker.User.Username, worker.Interval)

	respondSuccess(w)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad request method"), http.StatusBadRequest)
		return
	}

	var err error
	var worker shared.Worker

	if worker, err = getWorkerFromRequest(r); err != nil {
		httperr(w, err, http.StatusBadRequest)
		return
	}
	if W[worker.ID].Running == 1 {
		httperr(w, fmt.Errorf("worker is already running"), http.StatusBadRequest)
		return
	}
	worker.Running = 1
	if err = shared.Db.UpdateWorker(worker); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(worker.ID); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	M.Lock()
	W[worker.ID] = worker
	M.Unlock()

	go Worker(worker)

	log.Printf("Worker #%d started: { city: %s; user: %s; interval: %ds }\n", worker.ID, worker.City.Name, worker.User.Username, worker.Interval)

	respondSuccess(w)
}

// Worker is a worker thread
func Worker(w shared.Worker) {
	M.RLock()
	id := w.ID
	city := w.City
	user := w.User
	interval := w.Interval
	M.RUnlock()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	log.Printf("Worker #%d started: { city: %s; user: %s; interval: %ds }\n", id, city.Name, user.Username, interval)

	var status shared.WeatherStatus
	var err error
	var statusid int64
	for {
		select {
		case <-ticker.C:
			{
				if status, err = weather.GetWeather(city.ID); err != nil {
					log.Printf("Worker #%d error: { %s }\n", id, err.Error())
					continue
				}
				status.WorkerID = id
				if statusid, err = shared.Db.AddWeatherStatus(status); err != nil {
					log.Printf("Worker #%d error: %s\n", id, err.Error())
					continue
				}

				log.Printf("Worker #%d write: { city: %s; user: %s; status: #%d; time: %d; temp: %f }\n", id, city.Name, user.Username, statusid, status.Timestamp, status.Temperature)
			}
		case <-W[id].Stop:
			w = W[id]
			W[id] = w
			return
		}
	}
}

// StartService starts web service
func StartService() {
	var err error
	if W, err = shared.Db.GetWorkers(); err != nil {
		panic(err)
	}

	for _, w := range W {
		if w.Running == 1 {
			go Worker(w)
		} else {
			log.Printf("Worker #%d in paused\n", w.ID)
		}
	}

	log.Printf("Starting worker http service\n")
	mux := http.NewServeMux()

	mux.HandleFunc("/worker/add", addHandler)
	mux.HandleFunc("/worker/delete", deleteHandler)
	mux.HandleFunc("/worker/edit", editHandler)
	mux.HandleFunc("/worker/pause", pauseHandler)
	mux.HandleFunc("/worker/start", startHandler)
	mux.HandleFunc("/health-check", healthCheck)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", shared.WorkerPort), mux))
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
