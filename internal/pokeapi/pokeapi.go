package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thrashdev/bootdev-pokedex/internal/pokecache"
)

const locationURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"
const cacheInterval = time.Duration(10 * time.Second)

type Config struct {
	Previous *string
	Next     string
	Cache    *pokecache.Cache
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

type LocationDetailsResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokemonResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewConfig() *Config {
	conf := &Config{
		Next:     locationURL,
		Previous: nil,
		Cache:    pokecache.NewCache(cacheInterval),
	}
	return conf
}

func (config *Config) getLocations(url string) ([]Location, error) {
	body, ok := config.Cache.Get(url)
	if ok {

	}
	res, err := http.Get(url)
	if err != nil {
		return []Location{}, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return []Location{}, err
	}

	config.Cache.Add(url, body)

	locRes := LocationsResponse{}
	err = json.Unmarshal(body, &locRes)
	if err != nil {
		return []Location{}, err
	}
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

func (config *Config) getLocationsFromCache(cached []byte) ([]Location, error) {
	locRes := LocationsResponse{}
	err := json.Unmarshal(cached, &locRes)
	if err != nil {
		return []Location{}, err
	}
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

func GetNextLocations(config *Config) ([]Location, error) {
	cached, ok := config.Cache.Get(config.Next)
	if !ok {
		return config.getLocations(config.Next)
	}
	fmt.Println("READING FROM CACHE")
	return config.getLocationsFromCache(cached)

}

func GetPreviousLocations(config *Config) ([]Location, error) {
	if config.Previous == nil {
		return []Location{}, fmt.Errorf("Can't get previous locations, make a call to next locations first")
	}
	cached, ok := config.Cache.Get(*config.Previous)
	if !ok {
		return config.getLocations(*config.Previous)
	}
	fmt.Println("READING FROM CACHE")
	return config.getLocationsFromCache(cached)
}

func getPokemon(config *Config, url string) ([]Pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []Pokemon{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Pokemon{}, err
	}

	config.Cache.Add(url, body)

	var locDetails LocationDetailsResponse
	err = json.Unmarshal(body, &locDetails)
	if err != nil {
		return []Pokemon{}, err
	}

	result := []Pokemon{}
	for _, item := range locDetails.PokemonEncounters {
		result = append(result, item.Pokemon)
	}

	return result, nil

}

func getPokemonFromCache(cached []byte) ([]Pokemon, error) {
	var locDetails LocationDetailsResponse
	err := json.Unmarshal(cached, &locDetails)
	if err != nil {
		return []Pokemon{}, err
	}

	result := []Pokemon{}
	for _, item := range locDetails.PokemonEncounters {
		result = append(result, item.Pokemon)
	}

	return result, nil
}

func GetPokemon(config *Config, locName string) ([]Pokemon, error) {
	if locName == "" {
		return []Pokemon{}, fmt.Errorf("Please provide a name of the area to explore")
	}
	url := locationURL + "/" + locName + "/"

	cached, ok := config.Cache.Get(url)
	if !ok {
		return getPokemon(config, url)
	}
	fmt.Println("READING FROM CACHE")
	return getPokemonFromCache(cached)
}

func GetPokemonByName(config *Config, name string) (Pokemon, error) {
	if name == "" {
		return Pokemon{}, fmt.Errorf("Please provide a name of the Pokemon")
	}
	url := pokemonURL + "/" + name + "/"

}
