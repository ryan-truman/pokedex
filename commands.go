package main

import (
	"fmt"
	"os"
	"pokedex/internal/pokeapi"
)

type config struct {
	cliCommands map[string]cliCommand
	Location    map[string]string
	client      pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

func commandExit(*config, string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil 
}

func commandHelp(*config, string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}


func commandMap( sharedState *config, argument string) error {
	var url string
	if sharedState.Location["Next"] == "" {
		url = ""
	} else {
		url = sharedState.Location["Next"]
	}
	locations, err := sharedState.client.ListLocations(url)
	if err != nil {
		return err
	}
	sharedState.Location["Next"] = locations.Next
	sharedState.Location["Previous"] = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB( sharedState *config, argument string) error {
	var url string
	if sharedState.Location["Previous"] == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = sharedState.Location["Previous"]
	}
	locations, err := sharedState.client.ListLocations(url)
	if err != nil {
		return err
	}
	sharedState.Location["Next"] = locations.Next
	sharedState.Location["Previous"] = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}


func commandExplore( sharedState *config, area string) error {
	err := sharedState.client.ExploreArea(area)
	if err != nil {
		return err
	}
	return nil
}
