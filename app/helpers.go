package main

import (
	"os"
	"sort"
	"strings"
)

func Contains(arr []string, target string) int {
	for i, v := range arr {
		if target == v {
			return i
		}
	}
	return -1
}

func GetCmdPath(cmd string) string {
	path_env, exists := os.LookupEnv("PATH")
	if !exists {
		return ""
	}
	paths := strings.Split(path_env, ":")
	for _, path := range paths {
		if CheckIfPathExists(path + "/" + cmd) {
			return path + "/" + cmd
		}
	}
	return ""
}

func CheckIfPathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func CreateFile(path string) {
	if !CheckIfPathExists(path) {
		file, _ := os.Create(path)
		defer file.Close()
	}
}

func WriteToFile(path string, content string) {
	file, _ := os.OpenFile(path, os.O_WRONLY, 0644)
	file.WriteString(content)
	defer file.Close()
}

func AppendToFile(path string, content string) {
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(content)
}

func AutoComplete(input string) []string {
	builtins := [...]string{"echo", "type", "exit", "pwd"}
	suggestions := make(map[string]int)
	for _, cmd := range builtins {
		if strings.HasPrefix(cmd, input) {
			// suggestions = append(suggestions, cmd)
			suggestions[cmd] = 1
		}
	}

	path_env, exists := os.LookupEnv("PATH")
	if !exists {
		return nil
	}
	paths := strings.Split(path_env, ":")
	for _, path := range paths {
		entries, _ := os.ReadDir(path)
		for _, item := range entries {
			name := item.Name()
			if !item.IsDir() && strings.HasPrefix(name, input) {
				// suggestions = append(suggestions, name)
				suggestions[name] = 1
			}
		}
	}

	var suggestionsList []string
	for suggestion := range suggestions {
		suggestionsList = append(suggestionsList, suggestion)
	}

	sort.Strings(suggestionsList)
	return suggestionsList
}
