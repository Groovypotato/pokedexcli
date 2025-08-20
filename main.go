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
	callback func(*config,  ...string) error
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

func commandMap(cfg *config , args ...string) error {
	getURL := cfg.baseURL
	if cfg.nextUrl != "" {
		getURL = cfg.nextUrl
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

	cfg.nextUrl = locationAreas.Next
	cfg.prevUrl = locationAreas.Previous

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevUrl == ""{
		return errors.New("you're on the first page")
	}

	res, err := http.Get(cfg.prevUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("the call to location-area did not succeed: %v", res.StatusCode)
	}
	var locationareas LocationAreas
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&locationareas)
	if jsonErr != nil {
		return err
	}

	cfg.nextUrl = locationareas.Next
	cfg.prevUrl = locationareas.Previous

	for _, location := range locationareas.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandExplore( cfg *config, args ...string) error {
	 if len(args) == 0 {
        return errors.New("no location area provided")
    }
    area := args[0]
    fullURL := cfg.baseURL+area
    res, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("the call to location-area did not succeed: %v", res.StatusCode)
	}
	var locationarea LocationArea
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&locationarea)
	if jsonErr != nil {
		return err
	}


	for _, location := range locationarea.PokemonEncounters {
		fmt.Println(location.Pokemon.Name)
	}
	
	return nil
}




func commandHelp(cfg *config, args ...string) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n\n")
	for _, command := range commands{
		fmt.Printf("%s: %s\n",command.name,command.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
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
	var currentCfg config 
	currentCfg.baseURL = "https://pokeapi.co/api/v2/location-area/"
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		cmdName := userInput[0]
		if cmdFunc, ok := commands[cmdName]; ok {
			args := userInput[1:]
			err := cmdFunc.callback(&currentCfg,args...)
			if err != nil{
				fmt.Print(err)
			}
			} else {
			fmt.Print("Unknown command\n")
		}
	}
}