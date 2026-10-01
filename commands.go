package main

import (
	"fmt"
	"os"
	"math/rand/v2"
	"pokedex/internal/pokeapi"
)

type config struct {
	cliCommands map[string]cliCommand
	Location    map[string]string
	client      pokeapi.Client
	pokeDex     map[string]pokeapi.Pokemon
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


func commandCatch( sharedState *config, pokemon string) error {
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)
	targetPokemon, err := sharedState.client.GetPokemon(pokemon)
	if err != nil {
		return err
	}

	if rand.IntN(650) < min(600, targetPokemon.BaseExperience) {
		fmt.Printf("%v escaped!\n", targetPokemon.Name)
	} else {
		sharedState.pokeDex[targetPokemon.Name] = targetPokemon
		fmt.Printf("%v was caught!\n", targetPokemon.Name)
	}
	return nil
}

func commandInspect( sharedState *config, pokemon string) error {
	targetPokemon, ok := sharedState.pokeDex[pokemon]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %v\n", targetPokemon.Name)
	fmt.Printf("Height: %v\n", targetPokemon.Height)
	fmt.Printf("Weight: %v\n", targetPokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range targetPokemon.Stats{
		fmt.Printf("    -%v:%v \n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, Each := range targetPokemon.Types{
		fmt.Printf("    - %v\n", Each.Type.Name)
	}
	return nil
}

func commandPokedex( sharedState *config, args string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range sharedState.pokeDex {
		fmt.Printf("- %v\n", pokemon.Name)
	}
	return nil
}
