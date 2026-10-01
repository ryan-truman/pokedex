package pokeapi

import (
	"net/http"
	"io"
	"time"
	"fmt"
	"encoding/json"
	"pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

type locationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type areaPokemon struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func NewClient(timeout, interval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(interval),
	}
}

func (c *Client) ListLocations(url string) (locationAreas, error){
	var locations locationAreas
	if url == "" {
		url = baseURL + "/location-area"
	}

	body, cached := c.cache.Get(url)
	if !cached {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return locations, err
		}

		response, err := c.httpClient.Do(req)
		if err != nil {
			return locations, err
		}
		defer response.Body.Close()

		if response.StatusCode > 299 {
			return locations, fmt.Errorf("bad status code: %v", response.StatusCode)
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return locations, err
		}
		c.cache.Add(url, body)
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, err
	}

	return locations, nil
}

func (c *Client) ExploreArea(area string) error {
	fullURL := baseURL + "/location-area/" + area
	var pokemon areaPokemon
	body, cached := c.cache.Get(area)
	if !cached {
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return err
		}

		response, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode > 299 {
				return fmt.Errorf("bad status code: %v", response.StatusCode)
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		c.cache.Add(area, body)
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}

	for _, mon := range pokemon.PokemonEncounters {
		fmt.Println(mon.Pokemon.Name)
	}

	return nil
}


func (c *Client) GetPokemon(targetPokemon string) (Pokemon, error) {
	fullURL := baseURL + "/pokemon/" + targetPokemon
	var pokemon Pokemon
	body, cached := c.cache.Get(targetPokemon)
	if !cached {
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return pokemon, err
		}

		response, err := c.httpClient.Do(req)
		if err != nil {
			return pokemon, err
		}
		defer response.Body.Close()

		if response.StatusCode > 299 {
			return pokemon, fmt.Errorf("bad status code: %v", response.StatusCode)
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return pokemon, err
		}
		c.cache.Add(targetPokemon, body)
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
			return pokemon, err
	}
	return pokemon, nil
}
