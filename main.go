package main

import (
    "fmt"
    "bufio"
    "os"
)


type cliCommand struct {
    name string
    description string
    callback func() error
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
}
}

func commandHelp() error {
    fmt.Println("Type pokemon name and I will provide a description")
    return nil
}

func commandExit() error {
    fmt.Println("Exiting")
    return nil
}

func main() {

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("pokedex >")

        scanner.Scan()
        command := scanner.Text()

        cmd, ok := commands[command]
        if !ok {
            fmt.Println("Command does not exist: ", command) 
            continue
        }

        err := cmd.callback()
        if err != nil {
            fmt.Println("Error executing command: ", err)
            continue
        }

        if command == "exit" {
            return
        }
    }
}
