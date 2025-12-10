package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print
func main() {
	builtinCommands := make(map[string]struct{})
	builtinCommands["echo"] = struct{}{}
	builtinCommands["type"] = struct{}{}
	builtinCommands["exit"] = struct{}{}

	for {
		fmt.Print("$ ")
		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if command == "exit\n" {
			os.Exit(0)
		} else if command == "echo\n" {
			fmt.Println()
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Print(command[5:len(command)-1]+"\n")
		} else if command == "type\n" {
		} else if strings.HasPrefix(command, "type ") {
			com := strings.TrimSpace(command[5:len(command)-1])
			if _, exists := builtinCommands[com]; exists {
				fmt.Printf("%s is a shell builtin\n", com)
			} else {
				fmt.Println(com + ": not found")
			}
		} else{
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}