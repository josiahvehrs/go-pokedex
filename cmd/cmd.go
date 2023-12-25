package cmd

import (
	"fmt"
	"math/rand"
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
	Pokedex  map[string]poke.Pokemon
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
		"catch": {
			Name:        "catch",
			Description: "Catch a Pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a Pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "View your Pokedex",
			Callback:    commandPokedex,
		},
	}
}

func New() (map[string]Command, *Config) {
	commands := getCommands()
	c := cache.NewCache(5 * time.Minute)
	config := Config{Previous: "", Next: "", Cache: c, Pokedex: map[string]poke.Pokemon{}}
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
		return fmt.Errorf("explore command requires a location name")
	}

	detail, err := poke.GetLocationAreaDetail("https://pokeapi.co/api/v2/location-area/"+args[0], config.Cache)
	if err != nil {
		return fmt.Errorf("encountered an error fetching location area detail %s", err)
	}

	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")

	for _, p := range detail.PokemonEncounters {
		fmt.Println(" -", p.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch command requires a pokemon name")
	}

	pokemon, err := poke.GetPokemon("https://pokeapi.co/api/v2/pokemon/"+args[0], config.Cache)
	if err != nil {
		return fmt.Errorf("encountered an error fetching pokemon details %s", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	diceRoll := rand.Intn(400)
	time.Sleep(500 * time.Millisecond)
	if pokemon.BaseExperience > diceRoll {
		fmt.Println(pokemon.Name, "got away!")
		return nil
	}

	fmt.Println(pokemon.Name, "was caught!")
	fmt.Println("You may now inspect it.")
	config.Pokedex[pokemon.Name] = pokemon

	return nil
}

func commandInspect(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("inspect command requires a pokemon name")
	}
	pokemon, ok := config.Pokedex[args[0]]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)

	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	fmt.Println("Abilities:")
	for _, a := range pokemon.Abilities {
		fmt.Printf(" -slot-%d: %s\n", a.Slot, a.Ability.Name)
	}

	return nil
}

func commandPokedex(config *Config, args ...string) error {
	if len(config.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemon to get started!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, p := range config.Pokedex {
		fmt.Println(" -", p.Name)
	}

	return nil
}
