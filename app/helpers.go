package main

import (
	"os"
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
