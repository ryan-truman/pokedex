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
