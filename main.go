package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/thrashdev/bootdev-pokedex/internal/pokeapi"
)

type commandHandler func(*pokeapi.Config) error

type cliCommand struct {
	name        string
	description string
	handler     commandHandler
}

func handleCommand(cmdInput string, config *pokeapi.Config) error {
	defer fmt.Print("Pokedex > ")
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			handler:     commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the command",
			handler:     commandExit,
		},
		"map": {
			name:        "map",
			description: "Moves across the world of Pokemon, 20 locations at at time",
			handler:     commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Moves across the map 20 locations backwards",
			handler:     commandMapb,
		},
		"debug": {
			name:        "debug",
			description: "Prints the state of the config",
			handler:     commandDebug,
		},
	}
	command, ok := commands[cmdInput]
	if !ok {
		fmt.Printf("Unrecognized command: %s\n", cmdInput)
		return nil
	}
	err := command.handler(config)
	if err != nil {
		log.Printf("Encountered error while executing the command, err: %v\n", err)
	}

	return nil
}

func commandHelp(config *pokeapi.Config) error {
	fmt.Println("Welcome to the Pokedex:")
	fmt.Println("Usage")
	return nil
}

func commandExit(config *pokeapi.Config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *pokeapi.Config) error {
	locations, err := pokeapi.GetNextLocations(config)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil

}

func commandMapb(config *pokeapi.Config) error {
	locations, err := pokeapi.GetPreviousLocations(config)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil

}

func commandDebug(config *pokeapi.Config) error {
	fmt.Println(config)
	return nil
}

func main() {
	fmt.Print("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)
	config := pokeapi.NewConfig()
	for scanner.Scan() {
		handleCommand(scanner.Text(), config)
	}

}
