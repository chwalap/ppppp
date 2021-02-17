package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"shared"
)

type workerContext struct {
	Title   string
	Cities  map[int64]shared.City
	Workers []shared.Worker
}

// GetAvailableCitiesForUser returns cities for which user can create worker
func GetAvailableCitiesForUser(user shared.User) (map[int64]shared.City, error) {
	workers, err := shared.Db.GetUserWorkers(user)
	cities := map[int64]shared.City{}
	for cityid := range Cities {
		found := false
		for _, v := range workers {
			if v.CityID == cityid {
				found = true
				break
			}
		}
		if !found {
			cities[cityid] = Cities[cityid]
		}
	}
	return cities, err
}

func getWorkerContext(user shared.User) (workerContext, error) {
	workers, err := shared.Db.GetUserWorkers(user)
	cities, err := GetAvailableCitiesForUser(user)
	var w []shared.Worker
	for _, v := range workers {
		w = append(w, v)
	}
	sort.Slice(w, func(i, j int) bool { return w[i].ID < w[j].ID })
	context := workerContext{
		Title:   "Workers",
		Cities:  cities,
		Workers: w,
	}
	return context, err
}

// WorkersHandler handles workers page
func WorkersHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	var t *template.Template
	var context workerContext
	var user shared.User
	var err error

	w.Header().Set("Content-Type", "text/html")
	if t, err = template.ParseFiles(
		"html/workers.html",
		"html/head.html",
		"html/menu.html",
		"html/add_worker.html"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}

	userid := SessionManager.Get(r.Context(), "userid").(int64)
	if user, err = shared.Db.GetUserByID(userid); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if context, err = getWorkerContext(user); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, context); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
}

// AddWorkerHandler handles adding worker
func AddWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	// POST
	var worker shared.Worker
	var userid, cityid, interval int64
	var err error

	userid = SessionManager.Get(r.Context(), "userid").(int64)
	fmt.Printf("userid: %d\n", userid)

	if cityid, err = strconv.ParseInt(r.PostFormValue("city"), 10, 0); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	if interval, err = strconv.ParseInt(r.PostFormValue("interval"), 10, 0); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}

	worker = shared.Worker{
		UserID:   userid,
		CityID:   cityid,
		Interval: int(interval),
		Running:  1,
	}

	if err = sendWorkerRequest(worker, "add"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/worker", http.StatusSeeOther)
	return
}

// DeleteWorkerHandler handles deleting worker
func DeleteWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	var worker shared.Worker
	var workerid int64
	var err error

	if workerid, err = getIDFromURL(r); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if err = sendWorkerRequest(worker, "delete"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/worker", http.StatusSeeOther)
	return
}

type editContext struct {
	Cities   map[int64]shared.City
	Title    string
	CityID   int64
	ID       int64
	Interval int
	Running  int
}

func getEditContext(worker shared.Worker) editContext {
	cities, _ := GetAvailableCitiesForUser(worker.User) // todo: add error handling
	cities[worker.CityID] = worker.City
	return editContext{
		Title:    "Edit Worker",
		Cities:   cities,
		CityID:   worker.CityID,
		ID:       worker.ID,
		Interval: worker.Interval,
		Running:  worker.Running,
	}
}

// EditWorkerHandler handles editing worker
func EditWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var err error
	var workerid, interval int64
	var worker shared.Worker
	var t *template.Template

	if r.Method == "POST" {
		if workerid, err = strconv.ParseInt(r.PostFormValue("id"), 10, 0); err != nil {
			shared.HTTPerr(w, err, http.StatusBadRequest)
			return
		}
		if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
			shared.HTTPerr(w, err, http.StatusInternalServerError)
			return
		}
		if worker.CityID, err = strconv.ParseInt(r.PostFormValue("city"), 10, 0); err != nil {
			shared.HTTPerr(w, err, http.StatusBadRequest)
			return
		}
		if interval, err = strconv.ParseInt(r.PostFormValue("interval"), 10, 0); err != nil {
			shared.HTTPerr(w, err, http.StatusBadRequest)
			return
		}
		worker.Interval = int(interval)

		if err = sendWorkerRequest(worker, "edit"); err != nil {
			shared.HTTPerr(w, err, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/worker", http.StatusSeeOther)
		return
	}

	// GET
	w.Header().Set("Content-Type", "text/html")
	if t, err = template.ParseFiles(
		"html/empty.html",
		"html/head.html",
		"html/menu.html",
		"html/edit_worker.html"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if workerid, err = getIDFromURL(r); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, getEditContext(worker)); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
}

// PauseWorkerHandler handles pausing worker
func PauseWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	var worker shared.Worker
	var workerid int64
	var err error

	if workerid, err = getIDFromURL(r); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if err = sendWorkerRequest(worker, "pause"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/worker", http.StatusSeeOther)
	return
}

// StartWorkerHandler handles pausing worker
func StartWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	var worker shared.Worker
	var workerid int64
	var err error

	if workerid, err = getIDFromURL(r); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if err = sendWorkerRequest(worker, "start"); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/worker", http.StatusSeeOther)
	return
}

func getIDFromURL(r *http.Request) (int64, error) {
	for k, v := range r.URL.Query() {
		switch k {
		case "id":
			return strconv.ParseInt(v[0], 10, 0)
		}
	}
	return -1, fmt.Errorf("id not found")
}

func sendWorkerRequest(worker shared.Worker, method string) error {
	var result shared.WorkerResponse
	var response *http.Response
	var request *http.Request
	var jsonStr, body []byte
	var err error

	if jsonStr, err = json.Marshal(worker); err != nil {
		return err
	}
	if request, err = http.NewRequest("POST", shared.WorkerEndpoint+method, bytes.NewBuffer(jsonStr)); err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if response, err = client.Do(request); err != nil {
		return err
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf(result.Err)
	}
	return nil
}
