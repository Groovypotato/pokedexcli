package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name string
	description string
	callback func(*config,  ...string) error
}

type config struct {
	locationURL string
	pokemonURL string
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
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []PokemonStat `json:"stats"`
	Types []PokemonType `json:"types"`
}

type PokemonStat struct {
	Stat Stat `json:"stat"`
	BaseStat int `json:"base_stat"`
}

type Stat struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type PokemonType struct {
	Type PType `json:"type"`
}

type PType struct {
	Name string `json:"name"`
}


var commands map[string]cliCommand
var myPokedex map[string]Pokemon

func cleanInput(text string) []string {
	lWords := strings.ToLower(text)
	words := strings.Fields(lWords)
    return words
}

func commandMap(cfg *config , args ...string) error {
	getURL := cfg.locationURL
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
    fullURL := cfg.locationURL+area
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
        return errors.New("no location area provided")
    }
    pokemonToCatch := args[0]
	fullURL := cfg.pokemonURL+pokemonToCatch
	fmt.Printf("Throwing a Pokeball at %s...\n",pokemonToCatch)

	res, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("the call to pokemon did not succeed: %v\n", res.StatusCode)
	}
	var pokemon Pokemon
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonErr := json.Unmarshal(jsonData,&pokemon)
	if jsonErr != nil {
		return err
	}
	rand.NewSource(time.Now().UnixNano())
	chance := float64(rand.Intn(pokemon.BaseExperience))
	if chance > float64(pokemon.BaseExperience)*0.4 {
		fmt.Printf("%s was caught!\n",pokemon.Name)
		myPokedex[pokemon.Name] = pokemon
	}else {
		fmt.Printf("%s escaped!\n",pokemon.Name)
	}
	return nil
}

func comandInspect(cfg *config, args ...string) error {
	pInspect := args[0]
	pFound, ok := myPokedex[pInspect]
	if !ok {
		fmt.Printf("You have not caught %s to inspect. Go out there and get one Poke Master!\n", pInspect)
	} else {
		fmt.Printf("Name: %s\n",pFound.Name)
		fmt.Printf("Height: %v\n", pFound.Height)
		fmt.Printf("Weight: %v\n", pFound.Weight)
		fmt.Println("Stats:")
		for _, stat := range pFound.Stats {
			fmt.Printf("  -%s: %v\n",stat.Stat.Name,stat.BaseStat)
		}
		fmt.Println("Types:")
		for _,pktype := range pFound.Types {
			fmt.Printf("  - %s\n",pktype.Type.Name)
		}
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
		"catch": {
			name: "catch",
			description: "Attempts to catch a pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inpect",
			description: "Attempts to catch a pokemon",
			callback: comandInspect,
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
	myPokedex = make(map[string]Pokemon)
	currentCfg.locationURL = "https://pokeapi.co/api/v2/location-area/"
	currentCfg.pokemonURL = "https://pokeapi.co/api/v2/pokemon/"
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nPokedex > ")
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