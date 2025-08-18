package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
	lWords := strings.ToLower(text)
	words := strings.Fields(lWords)
    return words
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n\n")
	for _, command := range commands{
		fmt.Printf("%s: %s\n",command.name,command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}




func main() {
	commands = map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		value, found := commands[userInput[0]]
		if found {
			err := value.callback()
			if err != nil {
				fmt.Printf("Callback has failed: %v",err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}
	}
}