package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thrashdev/bootdev-pokedex/internal/pokeapi"
)

type commandHandler func(*pokeapi.Config, []string) error

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
		"explore": {
			name:        "explore",
			description: "Shows all Pokemon in the selected area",
			handler:     commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catches a pokemon named by the user",
			handler:     commandCatch,
		},
		"debug": {
			name:        "debug",
			description: "Prints the state of the config",
			handler:     commandDebug,
		},
	}
	input := strings.Split(cmdInput, " ")
	userCmd := input[0]
	arguments := input[1:]
	command, ok := commands[userCmd]
	if !ok {
		fmt.Printf("Unrecognized command: %s\n", cmdInput)
		return nil
	}
	err := command.handler(config, arguments)
	if err != nil {
		log.Printf("Encountered error while executing the command, err: %v\n", err)
	}

	return nil
}

func commandHelp(config *pokeapi.Config, arguments []string) error {
	fmt.Println("Welcome to the Pokedex:")
	fmt.Println("Usage")
	return nil
}

func commandExit(config *pokeapi.Config, arguments []string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *pokeapi.Config, arguments []string) error {
	locations, err := pokeapi.GetNextLocations(config)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil

}

func commandMapb(config *pokeapi.Config, arguments []string) error {
	locations, err := pokeapi.GetPreviousLocations(config)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc.Name)
	}
	return nil

}

func commandExplore(config *pokeapi.Config, arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("Please provide an area to explore")
	}
	locName := arguments[0]
	fmt.Printf("Exploring %s\n", locName)
	pokemon, err := pokeapi.GetPokemon(config, locName)
	if err != nil {
		return err
	}

	fmt.Println("Found pokemon: ")
	for _, pok := range pokemon {
		fmt.Println(pok.Name)
	}
	return nil
}

func commandCatch(config *pokeapi.Config, arguments []string) error {

}

func commandDebug(config *pokeapi.Config, arguments []string) error {
	fmt.Println("Config: ", config)
	fmt.Println("Arguments: ", arguments)
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
