package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const locationURL = "https://pokeapi.co/api/v2/location-area/"

type Config struct {
	Previous *string
	Next     string
}

type LocationsResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewConfig() *Config {
	return &Config{Next: locationURL}
}

func GetNextLocations(config *Config) ([]Location, error) {
	res, err := http.Get(config.Next)
	if err != nil {
		return []Location{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Location{}, err
	}

	locRes := LocationsResponse{}
	err = json.Unmarshal(body, &locRes)
	if err != nil {
		return []Location{}, err
	}
	fmt.Println(locRes.Next)
	if locRes.Previous != nil {
		fmt.Println(*locRes.Previous)
	}
	locs := []Location{}
	for _, elem := range locRes.Results {
		locs = append(locs, Location{Name: elem.Name, URL: elem.URL})
	}

	config.Next = locRes.Next
	config.Previous = locRes.Previous

	return locs, nil
}

func GetPreviousLocations(config *Config) ([]Location, error) {
	if config.Previous == nil {
		return []Location{}, fmt.Errorf("Can't get previous locations, make a call to next locations first")
	}
	res, err := http.Get(*config.Previous)
	if err != nil {
		return []Location{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Location{}, err
	}

	locRes := LocationsResponse{}
	err = json.Unmarshal(body, &locRes)
	if err != nil {
		return []Location{}, err
	}
	locs := []Location{}
	for _, elem := range locRes.Results {
		locs = append(locs, Location{Name: elem.Name, URL: elem.URL})
	}

	config.Next = locRes.Next
	config.Previous = locRes.Previous

	return locs, nil
}
