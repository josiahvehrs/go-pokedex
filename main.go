package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/josiahvehrs/go-pokedex/cmd"
)

func main() {
	commands, config := cmd.New()

	for {
		fmt.Print("pokedex > ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Printf("Encountered an error: %v", err)
		}

		cmd, ok := commands[strings.TrimSpace(strings.ToLower(scanner.Text()))]
		if !ok {
			fmt.Println("Sorry, I couldn't understand that. Try again.")
			continue
		}

		err = cmd.Callback(config)
		if err != nil {
			fmt.Printf("Encountered an error: %v\n", err)
		}
	}
}
