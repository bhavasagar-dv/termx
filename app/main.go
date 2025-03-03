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
		stdout := ""
		stderr := ""

		for i, arg := range args {
			if arg == "1>" || arg == ">" {
				stdout = args[i+1]
				args = args[:i]
			} else if arg == "2>" {
				stderr = args[i+1]
				args = args[:i]
			}
		}

		builtins := [...]string{"echo", "type", "exit", "pwd"}

		cmd_output := ""
		cmd_err := ""

		switch cmd {
		case "echo":
			cmd_output = strings.Join(args, " ")
		case "type":
			cmd_path := GetCmdPath(args[0])
			if Contains(builtins[:], args[0]) >= 0 {
				cmd_output = args[0] + " is a shell builtin"
			} else if len(cmd_path) > 0 {
				cmd_output = args[0] + " is " + cmd_path
			} else {
				cmd_output = args[0] + ": not found"
			}
		case "pwd":
			dir, _ := os.Getwd()
			cmd_output = dir
		case "cd":
			path := args[0]
			if len(path) > 0 && strings.Contains(path, "~") {
				home, _ := os.UserHomeDir()
				path = strings.ReplaceAll(path, "~", home)
			}

			if CheckIfPathExists(path) {
				os.Chdir(path)
			} else {
				cmd_output = "cd: " + path + ": No such file or directory"
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
				program.Stdin = os.Stdin
				program.Stderr = os.Stderr
				program.Stdout = os.Stdout
				err := program.Run()
				if err != nil {
					panic(err)
				}
			} else {
				cmd_err = cmd + ": command not found"
			}
		}

		if len(stdout) > 0 {
			CreateFile(stdout)
			WriteToFile(stdout, cmd_output)
		} else {
			fmt.Println(cmd_output)
		}

		if len(stderr) > 0 {
			CreateFile(stderr)
			WriteToFile(stderr, cmd_err)
		} else {
			fmt.Println(cmd_err)
		}
	}
}
