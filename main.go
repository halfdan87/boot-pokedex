package main

import (
    "fmt"
    "bufio"
    "errors"
    "strings"
    "os"
    "github.com/halfdan87/boot-pokedex/pokeapi"
)


type cliCommand struct {
    name string
    description string
    callback func(arg string) error
}

var commands map[string]cliCommand

func init() {
commands = map[string]cliCommand {
    "help": {
        name: "help",
        description: "Display help",
        callback: commandHelp,
    },
    "exit": {
        name: "exit",
        description: "Exit the program",
        callback: commandExit,
    },
    "map": {
        name: "map",
        description: "List locations, 20 next ones",
        callback: commandMap,
    },
    "mapb": {
        name: "mapb",
        description: "List locations, 20 previous ones",
        callback: commandMapb,
    },
    "explore": {
        name: "explore",
        description: "List pokemones in a specific location, provide location name as param",
        callback: commandExplore,
    },
}
}

func commandHelp(arg string) error {
    fmt.Println("Type pokemon name and I will provide a description")
    return nil
}

func commandExit(arg string) error {
    fmt.Println("Exiting")
    return nil
}

var defaultUrl string = "https://pokeapi.co/api/v2/location"
var pager pokeapi.Pagination = pokeapi.Pagination{
    Next : &defaultUrl,
}


func commandMap(arg string) error {
    if pager.Next == nil {
        return errors.New("No more next locations")
    }
    locations, newPager, err := pokeapi.GetLocations(pager.Next)
    if err != nil {
        return err
    }

    pager = *newPager

    for _, loc := range locations {
        fmt.Println(loc)
    }

    return nil
}

func commandMapb(arg string) error {
    if pager.Prev == nil {
        return errors.New("No more previous locations")
    }

    locations, newPager, err := pokeapi.GetLocations(pager.Prev)
    if err != nil {
        return err
    }

    pager = *newPager

    for _, loc := range locations {
        fmt.Println(loc)
    }

    return nil
}

func commandExplore(arg string) error {
    pokemons, err := pokeapi.GetPokemons(arg)
    if err != nil {
        return err
    }

    for _, pok := range pokemons {
        fmt.Println(pok)
    }

    return nil
}

func main() {

    scanner := bufio.NewScanner(os.Stdin)


    for {
        fmt.Print("pokedex >")

        scanner.Scan()
        command := scanner.Text()

        parts := strings.Split(command, " ")
        var arg string = ""

        if len(parts) == 2 {
            command = parts[0]
            arg = parts[1]
        }

        cmd, ok := commands[command]
        if !ok {
            fmt.Println("Command does not exist: ", command) 
            continue
        }

        err := cmd.callback(arg)
        if err != nil {
            fmt.Println("Error executing command: ", err)
            continue
        }

        if command == "exit" {
            return
        }
    }
}

