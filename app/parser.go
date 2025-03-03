package main

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
