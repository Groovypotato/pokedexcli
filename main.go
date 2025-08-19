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
	callback func(*config, string) error
}

type config struct {
	baseURL string
	nextUrl string
	prevUrl string
}

type LocationAreas struct {
	Count int
	Next string
	Previous string
	Results []Location
}

type Location struct {
	Name string
	URL string
}

type LocationArea struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
	lWords := strings.ToLower(text)
	words := strings.Fields(lWords)
    return words
}

func commandMap(config *config, location string) error {
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
	var locationAreas LocationAreas
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&locationAreas)
	if jsonErr != nil {
		return err
	}

	config.nextUrl = locationAreas.Next
	config.prevUrl = locationAreas.Previous

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandMapb(config *config, location string) error {
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
	var locationresponse LocationAreas
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
	
	return ni
}

func commandExplore( config *config,location string) {
	fullURL := config.baseURL+location
}



func commandHelp(config *config,location string) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n\n")
	for _, command := range commands{
		fmt.Printf("%s: %s\n",command.name,command.description)
	}
	return nil
}

func commandExit(config *config, location string) error {
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
		"explore": {
			name: "explore",
			description: "Displays pokemon in location",
			callback: commandExplore,
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