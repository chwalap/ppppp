package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"shared"
)

var apiKey = "d3e666adb659c9af7272d98c9627435f"

// Service for getting weather status
type Service interface {
	GetWeather(cityid int64) (shared.WeatherStatus, error)
}

// WeatherService gets and returns weather
type WeatherService struct{}

// GetWeather gets weather from the internet
func (s WeatherService) GetWeather(cityid int64) (shared.WeatherStatus, error) {
	var status shared.WeatherStatus
	var data map[string]interface{}
	var resp *http.Response
	var body []byte
	var ok bool
	var url string
	var err error

	// request
	url = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=metric", cityid, apiKey)
	if resp, err = http.Get(url); err != nil {
		return status, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return status, err
	}

	// parse response
	if err = json.Unmarshal(body, &data); err != nil {
		return status, err
	}
	if data == nil || data["main"] == nil || data["main"].(map[string]interface{})["temp"] == nil {
		return status, fmt.Errorf("weather service error {\nbody: %s\ndata: %v\n}", string(body), data)
	}
	if status.Temperature, ok = data["main"].(map[string]interface{})["temp"].(float64); !ok {
		return status, fmt.Errorf("cannot parse openweathermap response")
	}
	status.Timestamp = time.Now().Unix()

	return status, err
}
