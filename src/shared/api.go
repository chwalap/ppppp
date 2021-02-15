package shared

import (
	"fmt"
)

// WorkerPort is a port to worker service
var WorkerPort = 5002

// WorkerEndpoint point to worker endpoint
var WorkerEndpoint = fmt.Sprintf("http://worker:%d/worker/", WorkerPort)

// WeatherPort is a port to weather service
var WeatherPort = 5001

// WeatherEndpoint point to weather endpoint
var WeatherEndpoint = fmt.Sprintf("http://weather:%d/weather", WeatherPort)

// WorkerRequest describes add worker request
type WorkerRequest struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"userid"`
	CityID   int64 `json:"cityid"`
	Interval int   `json:"interval"`
}

// WorkerResponse describes worker service response
type WorkerResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"err,omitempty"`
}

// WeatherRequest describes weather service request
type WeatherRequest struct {
	CityID int64 `json:"cityid"`
}

// WeatherResponse describes weather service response
type WeatherResponse struct {
	Status WeatherStatus `json:"status"`
	Err    string        `json:"err,omitempty"`
}
