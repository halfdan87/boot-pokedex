package main

import (
    "fmt"
    "bufio"
    "errors"
    "os"
    "github.com/halfdan87/boot-pokedex/pokeapi"
)


type cliCommand struct {
    name string
    description string
    callback func(pager *pokeapi.Pagination) error
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
}
}

func commandHelp(pager *pokeapi.Pagination) error {
    fmt.Println("Type pokemon name and I will provide a description")
    return nil
}

func commandExit(pager *pokeapi.Pagination) error {
    fmt.Println("Exiting")
    return nil
}

func commandMap(pager *pokeapi.Pagination) error {
    if pager.Next == nil {
        return errors.New("No more next locations")
    }
    locations, newPager, err := pokeapi.GetLocations(pager.Next)
    if err != nil {
        return err
    }

    *pager = *newPager

    for _, loc := range locations {
        fmt.Println(loc)
    }

    return nil
}

func commandMapb(pager *pokeapi.Pagination) error {
    if pager.Prev == nil {
        return errors.New("No more previous locations")
    }

    locations, newPager, err := pokeapi.GetLocations(pager.Prev)
    if err != nil {
        return err
    }

    *pager = *newPager

    for _, loc := range locations {
        fmt.Println(loc)
    }

    return nil
}

func main() {

    scanner := bufio.NewScanner(os.Stdin)

    defaultUrl := "https://pokeapi.co/api/v2/location"

    pager := pokeapi.Pagination{
        Next : &defaultUrl,
    }

    for {
        fmt.Print("pokedex >")

        scanner.Scan()
        command := scanner.Text()

        cmd, ok := commands[command]
        if !ok {
            fmt.Println("Command does not exist: ", command) 
            continue
        }

        err := cmd.callback(&pager)
        if err != nil {
            fmt.Println("Error executing command: ", err)
            continue
        }

        if command == "exit" {
            return
        }
    }
}

