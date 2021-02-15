package shared

import (
	"encoding/json"
	"io/ioutil"
)

type cityJSON struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Country string `json:"country"`
	Coord   coords `json:"coord"`
}

type coords struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

// LoadCities loads all openweathermap cities
func LoadCities() (map[int64]City, error) {
	var cities map[int64]City = map[int64]City{}
	var jsonCities []cityJSON
	var fileContent []byte
	var err error
	if fileContent, err = ioutil.ReadFile("/db/cities.json"); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(fileContent, &jsonCities); err != nil {
		return nil, err
	}
	for _, c := range jsonCities {
		cities[c.ID] = City{ID: c.ID, Name: c.Name}
	}
	return cities, nil
}
