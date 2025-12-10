package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
)

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
						executables[file.Name()] = path
					}
				}
			}
		}
	}
	builtinCommands[ECHO] = true
	builtinCommands[EXIT] = true
	builtinCommands[TYPE] = true
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
		command = strings.TrimSpace(command)
		splits := strings.Split(command, " ")
		command = splits[0]
		args := splits[1:]

		command = strings.Split(command, "\n")[0]
		switch {
		case command == EXIT:
			handleExit()
		case command == ECHO:
			handleEcho(args)
		case command == TYPE:
			handleType(args)
		case func() bool { _, ok := executables[command]; return ok }():
			handleExecutable(command, executables[command], args)
		default:
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}

func handleExit() {
	os.Exit(0)
}

func handleEcho(args []string) {
	if len(args) > 0 {
		fmt.Println(strings.Join(args, " "))
	} else {
		fmt.Println()
	}
}

func handleType(args []string) {
	for _, arg := range args {
		if len(arg) > 0 {
			if _, exists := builtinCommands[arg]; exists {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else if _, path := executables[arg]; path {
				fmt.Printf("%s is %s\n", arg, executables[arg]+"/"+arg)
			} else {
				fmt.Println(arg + ": not found")
			}
		}
	}
}

func handleExecutable(command string, path string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
	}
	fmt.Print(string(out))
}
