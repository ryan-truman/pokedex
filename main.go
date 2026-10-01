package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"time"
	"pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	pokeapiClient := pokeapi.NewClient(5 * time.Second, 5 * time.Second)
	sharedState := config {
		cliCommands: map[string]cliCommand {
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Show the help menu",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Display the names of the next 20 locations",
				callback:    commandMap,
			},
			"mapb": {
				name:        "map",
				description: "Display the names of the previous 20 locations",
				callback:    commandMapB,
			},
			"explore": {
				name:        "explore",
				description: "Explore a location for pokemon",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Attempt to catch a pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect a caught pokemon",
				callback:    commandInspect,
			},
		},
		Location: map[string]string {
			"Next": "",
			"Previous": "",
		},
		client: pokeapiClient,
		pokeDex: map[string]pokeapi.Pokemon {},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		if scanner.Scan() {
			input := cleanInput(scanner.Text())
			argument := ""
			if len(input) > 1 {
				argument = input[1]
			}
			if len(input) > 0 {
				val, ok := sharedState.cliCommands[input[0]]
				if ok {
					err := val.callback(&sharedState, argument)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		} else if scanner.Err() != nil {
			fmt.Println(scanner.Err())
		} else {
			break
		}
	}
}
