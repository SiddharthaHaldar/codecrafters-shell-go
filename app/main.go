package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var builtinCommands = make(map[string]bool)
var executables = make(map[string]string)
var PATHs []string

func init() {
	PATH := os.Getenv("PATH")
	PATHs = strings.Split(PATH, ":")
	for _, path := range PATHs {
		files, err := os.ReadDir(path)
		for _, file := range files {
			if !file.IsDir() || err != nil {
				info, _ := os.Stat(path + "/" + file.Name())
				if info != nil {
					perms := info.Mode().Perm().String()
					if strings.Contains(perms, "x") {
						executables[file.Name()] = path + "/" + file.Name()
					}
				}
			}
		}
	}
	builtinCommands["echo"] = true
	builtinCommands["type"] = true
	builtinCommands["exit"] = true
}

func main() {
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
			fmt.Print(command[5:len(command)-1] + "\n")
		} else if command == "type\n" {
		} else if strings.HasPrefix(command, "type ") {
			com := strings.TrimSpace(command[5 : len(command)-1])
			if _, exists := builtinCommands[com]; exists {
				fmt.Printf("%s is a shell builtin\n", com)
			} else if _, path := executables[com]; path {
				fmt.Printf("%s is %s\n", com, executables[com])
			} else {
				fmt.Println(com + ": not found")
			}
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
