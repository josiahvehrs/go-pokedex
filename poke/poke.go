package poke

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/josiahvehrs/go-pokedex/cache"
)

type AreaLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetAreaLocations(url string, c *cache.Cache) (AreaLocations, error) {
	var locations AreaLocations
	var body []byte

	body, ok := c.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return locations, err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return locations, err
		}

		c.Add(url, body)
	}

	err := json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}

	return locations, nil
}
