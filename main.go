package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func handleCommand(cmdInput string) error {
	defer fmt.Print("Pokedex > ")
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the command",
			callback:    commandExit,
		},
	}
	command, ok := commands[cmdInput]
	if ok {
		command.callback()
	} else {
		fmt.Printf("Unrecognized command: %s\n", cmdInput)
	}

	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex:")
	fmt.Println("Usage")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	fmt.Print("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleCommand(scanner.Text())
	}

}
