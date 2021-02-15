package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"shared"
	"weather/service"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	log.Printf("Weather service startup\n")

	weatherHandler := httptransport.NewServer(
		service.MakeWeatherEndpoint(service.WeatherService{}),
		decodeRequest,
		encodeResponse,
	)

	http.Handle("/weather", weatherHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", shared.WeatherPort), nil))
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	log.Printf("Encoding response: %v\n", response)
	return json.NewEncoder(w).Encode(response)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request shared.WeatherRequest
	var err error

	log.Printf("Decoding request: %v\n", r)
	err = json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
