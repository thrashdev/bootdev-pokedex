package pokeapi

import (
	"archive/tar"
	"encoding/json"
	"io"
	"net/http"
)

const locationURL = "https://pokeapi.co/api/v2/location-area/"

type ApiConfig struct {
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

func GetNextLocations(conf *ApiConfig) ([]Location, error) {
	targetURL := conf.Next
	if targetURL == "" {
		targetURL = locationURL
	}

	res, err := http.Get(targetURL)
	if err != nil {
		return []Location{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Location{}, err
	}

	locRes := LocationsResponse{}
	err = json.Unmarshal(body, locRes)
	if err != nil {
		return []Location{}, err
	}

	return nil
}
