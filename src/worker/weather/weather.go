package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"shared"
)

// GetWeather returns weather from another microservice
func GetWeather(cityid int64) (shared.WeatherStatus, error) {
	var result shared.WeatherResponse
	var request *http.Request
	var response *http.Response
	var status shared.WeatherStatus
	var jsonStr []byte
	var body []byte
	var err error

	// prepare request
	if jsonStr, err = json.Marshal(shared.WeatherRequest{CityID: cityid}); err != nil {
		return status, err
	}
	if request, err = http.NewRequest("POST", shared.WeatherEndpoint, bytes.NewBuffer(jsonStr)); err != nil {
		return status, err
	}
	request.Header.Set("Content-Type", "application/json")

	// query API
	client := &http.Client{}
	if response, err = client.Do(request); err != nil {
		return status, err
	}
	defer response.Body.Close()

	// read response
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return status, err
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return status, err
	}
	if result.Err != "" {
		return status, fmt.Errorf("%s", result.Err)
	}
	return result.Status, nil
}
