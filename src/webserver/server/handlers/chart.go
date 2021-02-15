package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"

	"shared"
)

const dateFormat = "2006-01-02"

func getURLParameters(r *http.Request) (shared.Worker, int64, error) {
	var worker shared.Worker
	var workerid int64
	var timestamp int64 = -1
	var date time.Time
	var err error

	for k, v := range r.URL.Query() {
		switch k {
		case "id":
			if workerid, err = strconv.ParseInt(v[0], 10, 0); err != nil {
				return worker, timestamp, err
			}
			if worker, err = shared.Db.GetWorkerByID(workerid); err != nil {
				return worker, timestamp, err
			}
		case "date":
			if date, err = time.Parse(dateFormat, v[0]); err != nil {
				return worker, timestamp, err
			}
			timestamp = date.Unix()
		}
	}

	return worker, timestamp, err
}

type chartContext struct {
	Title   string
	ID      int64
	City    string
	Stats   []shared.WeatherStatus
	Date    string
	MinDate string
	MaxDate string
}

func getChartContext(worker shared.Worker, timestamp int64) (chartContext, error) {
	var context chartContext
	var stats []shared.WeatherStatus
	var err error
	var dateStart, dateEnd int64

	context.Title = "Chart"
	context.ID = worker.ID
	context.City = worker.City.Name

	if dateStart, dateEnd, err = shared.Db.GetGatheringStartEnd(worker); err != nil {
		return context, err
	}
	if dateStart == 0 || dateEnd == 0 {
		return context, fmt.Errorf("cannot get start and end dates")
	}
	context.Date = time.Now().Format(dateFormat)
	context.MinDate = time.Unix(dateStart, 0).Format(dateFormat)
	context.MaxDate = time.Unix(dateEnd, 0).Format(dateFormat)

	if timestamp != -1 {
		dateStart = timestamp
		dateEnd = time.Unix(timestamp, 0).AddDate(0, 0, 1).Unix()
		context.Date = time.Unix(timestamp, 0).Format(dateFormat)
	}
	if stats, err = shared.Db.ReadWeatherStatuses(worker, dateStart, dateEnd); err != nil {
		return context, err
	}
	// remove redundant values
	context.Stats = append(context.Stats, stats[0])
	last := context.Stats[0]
	for i := 0; i < len(stats); i++ {
		if math.Abs(last.Temperature-stats[i].Temperature) > 0.001 {
			context.Stats = append(context.Stats, stats[i])
			last = context.Stats[len(context.Stats)-1]
		}
	}
	if len(stats) >= 2 {
		if len(context.Stats) >= 2 {
			context.Stats[len(context.Stats)-1] = stats[len(stats)-1]
		} else {
			context.Stats = append(context.Stats, stats[len(stats)-1])
		}
	}
	return context, err
}

// ChartHandler handles worker chart
func ChartHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		httperr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	var worker shared.Worker
	var context chartContext
	var timestamp int64
	var t *template.Template
	var err error

	if worker, timestamp, err = getURLParameters(r); err != nil {
		httperr(w, err, http.StatusBadRequest)
		return
	}
	userid := SessionManager.Get(r.Context(), "userid").(int64)
	if worker.UserID != userid {
		httperr(w, fmt.Errorf("worker does not belong to user"), http.StatusBadRequest)
		return
	}
	if context, err = getChartContext(worker, timestamp); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if t, err = template.ParseFiles(
		"html/chart.html",
		"html/scripts.html",
		"html/date.html",
		"html/menu.html",
		"html/head.html"); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, context); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
}
