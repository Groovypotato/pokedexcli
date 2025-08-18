package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lWords := strings.ToLower(text)
	words := strings.Fields(lWords)
    return words
}



func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		fmt.Printf("Your command was: %s\n",userInput[0])
	}
}