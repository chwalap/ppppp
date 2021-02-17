package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"shared"
)

var apiKey = "d3e666adb659c9af7272d98c9627435f"

// StartService starts web service
func StartService() {
	log.Printf("Starting weather web service\n")
	mux := http.NewServeMux()

	mux.HandleFunc("/weather", weatherHandler)
	mux.HandleFunc("/health-check", healthCheck)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", shared.WeatherPort), mux))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	var request shared.WeatherRequest
	var response shared.WeatherResponse
	var data map[string]interface{}
	var resp *http.Response
	var body []byte
	var ok bool
	var url string
	var err error

	// decode request
	if r.Method != "POST" {
		shared.HTTPerr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&r); err != nil {
		shared.HTTPerr(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("Recieved request: { %v }\n", request)

	// send weather request
	url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=metric", request.CityID, apiKey)
	if resp, err = http.Get(url); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}

	// parse weather response
	if err = json.Unmarshal(body, &data); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
	if data == nil || data["main"] == nil || data["main"].(map[string]interface{})["temp"] == nil {
		shared.HTTPerr(w, fmt.Errorf("weather service error {\nbody: %s\ndata: %v\n}", string(body), data), http.StatusInternalServerError)
		return
	}
	if response.Status.Temperature, ok = data["main"].(map[string]interface{})["temp"].(float64); !ok {
		shared.HTTPerr(w, fmt.Errorf("cannot parse openweathermap response"), http.StatusInternalServerError)
		return
	}
	response.Status.Timestamp = time.Now().Unix()
	response.Status.WorkerID = request.WorkerID
	log.Printf("Sent response: { %v }\n", response)

	// encode response
	if err = json.NewEncoder(w).Encode(response); err != nil {
		shared.HTTPerr(w, err, http.StatusInternalServerError)
		return
	}
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
