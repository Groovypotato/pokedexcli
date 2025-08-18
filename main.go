package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

type config struct {
	baseURL string
	nextUrl string
	prevUrl string
}

type LocationResponse struct {
	Count int
	Next string
	Previous string
	Results []Location
}

type Location struct {
	Name string
	URL string
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
	lWords := strings.ToLower(text)
	words := strings.Fields(lWords)
    return words
}

func commandMap(config *config) error {
	getURL := config.baseURL
	if config.nextUrl != "" {
		getURL = config.nextUrl
	}
	res, err := http.Get(getURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("the call to location-area did not succeed: %v", res.StatusCode)
	}
	var locationresponse LocationResponse
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&locationresponse)
	if jsonErr != nil {
		return err
	}

	config.nextUrl = locationresponse.Next
	config.prevUrl = locationresponse.Previous

	for _, location := range locationresponse.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandMapb(config *config) error {
	if config.prevUrl == ""{
		return errors.New("you're on the first page")
	}

	res, err := http.Get(config.prevUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("the call to location-area did not succeed: %v", res.StatusCode)
	}
	var locationresponse LocationResponse
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&locationresponse)
	if jsonErr != nil {
		return err
	}

	config.nextUrl = locationresponse.Next
	config.prevUrl = locationresponse.Previous

	for _, location := range locationresponse.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}



func commandHelp(config *config) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n\n")
	for _, command := range commands{
		fmt.Printf("%s: %s\n",command.name,command.description)
	}
	return nil
}

func commandExit(config *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}




func main() {
	commands = map[string]cliCommand{
		"map": {
			name: "map",
			description: "Displays 20 locations at a time",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays last 20 locations",
			callback: commandMapb,
		},
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
	var currentConfig config 
	currentConfig.baseURL = "https://pokeapi.co/api/v2/location-area/"
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		value, found := commands[userInput[0]]
		if found {
			err := value.callback(&currentConfig)
			if err != nil {
				fmt.Printf("Callback has failed: %v",err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}
	}
}