package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func ReadInput(stdin io.Reader) string {
	input := ""

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	reader := bufio.NewReader(stdin)

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
			if len(suffix) > 1 {
				input += suffix
				fmt.Printf("%s", suffix)
			} else {
				fmt.Printf("\a")
			}
		default:
			fmt.Printf("%s", char)
			input += char
		}
	}
	return input
}

func ExtractArgsAndCmd(input_str string) (string, []string) {
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
