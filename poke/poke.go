package poke

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/josiahvehrs/go-pokedex/cache"
)

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetail struct {
	ID       int `json:"id"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height         int `json:"height"`
	ID             int `json:"id"`
	Abilities      []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	Name  string `json:"name"`
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
	Weight int `json:"weight"`
}

func GetLocationAreas(url string, c *cache.Cache) (LocationArea, error) {
	var locations LocationArea
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

func GetLocationAreaDetail(url string, c *cache.Cache) (LocationAreaDetail, error) {
	var detail LocationAreaDetail
	var body []byte

	body, ok := c.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return detail, err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return detail, err
		}

		c.Add(url, body)
	}

	err := json.Unmarshal(body, &detail)
	if err != nil {
		return detail, err
	}

	return detail, nil
}

func GetPokemon(url string, c *cache.Cache) (Pokemon, error) {
	var pokemon Pokemon
	var body []byte

	body, ok := c.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return pokemon, err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return pokemon, err
		}

		c.Add(url, body)
	}

	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
