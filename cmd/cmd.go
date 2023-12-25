package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/josiahvehrs/go-pokedex/cache"
	"github.com/josiahvehrs/go-pokedex/poke"
)

type Command struct {
	Name        string
	Description string
	Callback    func(config *Config) error
}

type Config struct {
	Previous string
	Next     string
	Cache    *cache.Cache
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			Name:        "help",
			Description: "Prints the help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Get next 20 map locations",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Get previous 20 map locations",
			Callback:    commandMapBack,
		},
	}
}

func New() (map[string]Command, *Config) {
	commands := getCommands()
	c := cache.NewCache(5 * time.Minute)
	config := Config{Previous: "", Next: "", Cache: c}
	return commands, &config
}

func commandHelp(config *Config) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!", "Usage:")
	fmt.Println()
	for _, value := range commands {
		fmt.Printf("%s - %s\n", value.Name, value.Description)
	}
	fmt.Println()
	return nil
}

func commandExit(config *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *Config) error {
	var url string
	if config.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	} else {
		url = config.Next
	}

	locations, err := poke.GetLocationAreas(url, config.Cache)
	if err != nil {
		return fmt.Errorf("encountered an error fetching locations %s", err)
	}

	config.Previous = locations.Previous
	config.Next = locations.Next

	for _, area := range locations.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}

func commandMapBack(config *Config) error {
	if config.Previous == "" {
		return fmt.Errorf("no previous locations")
	}

	locations, err := poke.GetLocationAreas(config.Previous, config.Cache)
	if err != nil {
		return fmt.Errorf("encountered an error fetching locations %s", err)
	}

	config.Previous = locations.Previous
	config.Next = locations.Next

	for _, area := range locations.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}
