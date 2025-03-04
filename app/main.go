package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"
)

func main() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)
		input := ""

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

		if err != nil {

			panic(err)

		}

		defer term.Restore(int(os.Stdin.Fd()), oldState)

	loop:
		for {
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading the input: ", err)
				os.Exit(1)
			}

			char := string(b)
			switch b {
			case '\x03':
				// ctrl + c
				os.Exit(0)
			case '\x7F':
				// backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
				fmt.Printf("\b \b")
			case '\n', '\r':
				// new line
				fmt.Printf("\n")
				break loop
			case '\t':
				suggestion := AutoComplete(input)
				suffix := suggestion[len(input):] + " "
				input += suffix
				fmt.Printf("%s", suffix)
			default:
				fmt.Printf("%s", char)
				input += char
			}
		}

		cmd, args := ExtractArgsAndCmd(strings.TrimSpace(input))
		stdout := ""
		stderr := ""
		append := false

		for i, arg := range args {
			if arg == "1>" || arg == ">" || arg == ">>" || arg == "1>>" {
				stdout = args[i+1]
				args = args[:i]
			} else if arg == "2>" || arg == "2>>" {
				stderr = args[i+1]
				args = args[:i]
			}

			if arg == ">>" || arg == "1>>" || arg == "2>>" {
				append = true
			}
		}

		builtins := [...]string{"echo", "type", "exit", "pwd"}

		cmd_output := ""
		cmd_err := ""

		switch cmd {
		case "echo":
			cmd_output = strings.Join(args, " ") + "\n"
		case "type":
			cmd_path := GetCmdPath(args[0])
			if Contains(builtins[:], args[0]) >= 0 {
				cmd_output = args[0] + " is a shell builtin"
			} else if len(cmd_path) > 0 {
				cmd_output = args[0] + " is " + cmd_path
			} else {
				cmd_output = args[0] + ": not found"
			}
			cmd_output += "\n"
		case "pwd":
			dir, _ := os.Getwd()
			cmd_output = dir + "\n"
		case "cd":
			path := args[0]
			if len(path) > 0 && strings.Contains(path, "~") {
				home, _ := os.UserHomeDir()
				path = strings.ReplaceAll(path, "~", home)
			}

			if CheckIfPathExists(path) {
				os.Chdir(path)
			} else {
				cmd_output = "cd: " + path + ": No such file or directory\n"
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

				var outErrBuffer bytes.Buffer
				var outBuffer bytes.Buffer
				program.Stderr = &outErrBuffer
				program.Stdout = &outBuffer

				program.Run()
				cmd_err = outErrBuffer.String()
				cmd_output = outBuffer.String()
			} else {
				cmd_err = cmd + ": command not found\n"
			}
		}

		if len(stdout) > 0 {
			CreateFile(stdout)
			if append {
				AppendToFile(stdout, cmd_output)
			} else {
				WriteToFile(stdout, cmd_output)
			}
		} else if len(cmd_output) > 0 {
			fmt.Print(cmd_output)
		}

		if len(stderr) > 0 {
			CreateFile(stderr)
			if append {
				AppendToFile(stderr, cmd_err)
			} else {
				WriteToFile(stderr, cmd_err)
			}
		} else if len(cmd_err) > 0 {
			fmt.Print(cmd_err)
		}
	}
}
