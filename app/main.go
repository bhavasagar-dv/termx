package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading the input: ", err)
			os.Exit(1)
		}

		runtime := runtime.GOOS
		max_cap := len(input) - 1
		if runtime == "windows" {
			max_cap--
		}

		cmd, args := ExtractArgsAndCmd(input[:max_cap])

		builtins := [...]string{"echo", "type", "exit", "pwd"}

		switch cmd {
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			cmd_path := GetCmdPath(args[0])
			if Contains(builtins[:], args[0]) >= 0 {
				fmt.Println(args[0], "is a shell builtin")
			} else if len(cmd_path) > 0 {
				fmt.Println(args[0], "is", cmd_path)
			} else {
				fmt.Println(args[0] + ": not found")
			}
		case "pwd":
			dir, _ := os.Getwd()
			fmt.Println(dir)
		case "cd":
			path := args[0]
			if len(path) > 0 && strings.Contains(path, "~") {
				home, _ := os.UserHomeDir()
				path = strings.ReplaceAll(path, "~", home)
			}

			if CheckIfPathExists(path) {
				os.Chdir(path)
			} else {
				fmt.Println("cd: " + path + ": No such file or directory")
			}

		case "exit":
			exit_status := 0
			if len(args) > 0 {
				exit_status, _ = strconv.Atoi(args[0])
			}
			os.Exit(exit_status)
		default:
			cmd_path := GetCmdPath(cmd)
			if len(cmd_path) > 0 {
				program := exec.Command(cmd, args...)
				program.Stderr = os.Stderr
				program.Stdout = os.Stdout
				err := program.Run()
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println(cmd + ": command not found")
			}
		}
	}
}
