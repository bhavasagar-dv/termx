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

		cmd, args := extractArgsAndCmd(input[:max_cap])

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
				// program.Stderr = os.Stderr
				program.Stdout = os.Stdout
				_, err := program.Output()
				if err != nil {
					panic(err)
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
		if checkIfPathExists(path + "/" + cmd) {
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

func extractArgsAndCmd(input_str string) (string, []string) {
	cmd := ""
	var args []string
	curr := ""
	open_single_quote := false
	open_double_quote := false
	prev_backslash := false
	for _, char := range input_str {
		if char == rune(' ') && cmd == "" && !open_single_quote && !open_double_quote && !prev_backslash {
			cmd = curr
			curr = ""
		} else if char == rune(' ') && !open_single_quote && !open_double_quote && !prev_backslash {
			if len(curr) > 0 {
				args = append(args, curr)
			}
			curr = ""
		}

		if char == rune('"') && !open_single_quote && !prev_backslash {
			open_double_quote = !open_double_quote
		} else if char == rune('\'') && !open_double_quote && !prev_backslash {
			open_single_quote = !open_single_quote
		} else if char == rune('\\') && !open_single_quote && !prev_backslash {
			prev_backslash = true
		} else if char == rune(' ') && !prev_backslash {
			if open_single_quote || open_double_quote {
				curr += string(char)
			}
		} else {
			if prev_backslash && open_double_quote && char != rune('\\') && char != rune('$') && char != rune('"') {
				curr += string('\\')
			}
			prev_backslash = false
			curr += string(char)
		}
	}
	if len(curr) > 0 {
		if cmd == "" {
			cmd = curr
		} else {
			args = append(args, curr)
		}
	}

	return cmd, args
}
