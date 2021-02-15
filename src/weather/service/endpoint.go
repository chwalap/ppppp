package service

import (
	"context"

	"shared"

	"github.com/go-kit/kit/endpoint"
)

// MakeWeatherEndpoint creates weather endpoint
func MakeWeatherEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		status, err := svc.GetWeather(request.(shared.WeatherRequest).CityID)
		if err != nil {
			return shared.WeatherResponse{Status: status, Err: err.Error()}, nil
		}
		return shared.WeatherResponse{Status: status, Err: ""}, nil
	}
}
