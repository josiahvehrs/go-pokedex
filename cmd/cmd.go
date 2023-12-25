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
	Callback    func(config *Config, args ...string) error
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
		"explore": {
			Name:        "explore",
			Description: "Explore a location",
			Callback:    commandExplore,
		},
	}
}

func New() (map[string]Command, *Config) {
	commands := getCommands()
	c := cache.NewCache(5 * time.Minute)
	config := Config{Previous: "", Next: "", Cache: c}
	return commands, &config
}

func commandHelp(config *Config, args ...string) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!", "Usage:")
	fmt.Println()
	for _, value := range commands {
		fmt.Printf("%s - %s\n", value.Name, value.Description)
	}
	fmt.Println()
	return nil
}

func commandExit(config *Config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *Config, args ...string) error {
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

func commandMapBack(config *Config, args ...string) error {
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

func commandExplore(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore command requires a location argument")
	}

	detail, err := poke.GetLocationAreaDetail("https://pokeapi.co/api/v2/location-area/"+args[0], config.Cache)
	if err != nil {
		return fmt.Errorf("encountered an error fetching location area detail %s", err)
	}

	fmt.Println("Exploring ", args[0], "...")
	fmt.Println("Found Pokemon:")

	for _, p := range detail.PokemonEncounters {
		fmt.Println(" - ", p.Pokemon.Name)
	}
	return nil
}
