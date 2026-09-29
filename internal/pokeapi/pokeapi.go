package pokeapi

import (
	"net/http"
	"io"
	"time"
	"encoding/json"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
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

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ListLocations(url string) (locationAreas, error){
	var locations locationAreas
	if url == "" {
		url = baseURL + "/location-area"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return locations, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return locations, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		return locations, err
	}
	if err != nil {
		return locations, err
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, err
	}

	return locations, nil
}
