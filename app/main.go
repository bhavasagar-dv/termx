package main

import (
	"bufio"
	"errors"
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

		input_eval := strings.Split(input[:max_cap], " ")
		if len(input_eval) == 0 {
			continue
		}

		cmd := input_eval[0]
		raw_args := input_eval[1:]
		var args []string
		for _, arg := range raw_args {
			arg = strings.Trim(arg, "\"' ")
			if len(arg) > 0 {
				args = append(args, arg)
			}
		}
		builtins := [...]string{"echo", "type", "exit", "pwd"}

		switch cmd {
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			cmd_path := getCmdPath(args[0])
			if contains(builtins[:], args[0]) >= 0 {
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

			if checkIfPathExists(path) {
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
			cmd_path := getCmdPath(cmd)
			if len(cmd_path) > 0 {
				program := exec.Command(cmd, args...)
				program.Stderr = os.Stderr
				program.Stdout = os.Stdout
				err := program.Run()
				if err != nil {
					fmt.Println(cmd + ": command not found")
				}
			} else {
				fmt.Println(cmd + ": command not found")
			}
		}
	}
}

func contains(arr []string, target string) int {
	for i, v := range arr {
		if target == v {
			return i
		}
	}
	return -1
}

func getCmdPath(cmd string) string {
	path_env, exists := os.LookupEnv("PATH")
	if !exists {
		return ""
	}
	paths := strings.Split(path_env, ":")
	for _, path := range paths {
		if _, err := os.Stat(path + "/" + cmd); !errors.Is(err, os.ErrNotExist) {
			return path + "/" + cmd
		}
	}
	return ""
}

func checkIfPathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
