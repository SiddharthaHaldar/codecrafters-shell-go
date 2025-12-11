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
	PWD  = "pwd"
	CD = "cd"
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
	builtinCommands[PWD] = true
	builtinCommands[CD] = true
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
			case func() bool { _, ok := builtinCommands[command]; return ok }():
				handleBuiltIn(command, args)
			case func() bool { _, ok := executables[command]; return ok }():
				handleExecutable(command, executables[command], args)
			default:
				fmt.Println(command + ": command not found")
		}
	}
}
func handleBuiltIn(command string, args []string) {
	switch command {
		case EXIT:
			handleExit()
		case ECHO:
			handleEcho(args)
		case TYPE:
			handleType(args)
		case PWD:
			handlePresentWorkingDirectory(args)
		case CD:
			handleChangeDirectory(args)
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

func handlePresentWorkingDirectory(args []string) {
	if len(args) > 0 {
		fmt.Println("pwd: too many arguments")
		return
	}	
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("pwd: error retrieving current directory")
		return
	}
	fmt.Println(pwd)
}

func handleChangeDirectory(args []string) {
	if len(args) == 0 {
		return
	} else {
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Printf("cd: %s: No such file or directory\n", args[0])
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
