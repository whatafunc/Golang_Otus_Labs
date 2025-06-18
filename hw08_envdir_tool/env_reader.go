package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Each file in the directory represents a variable where:
// - Filename is the variable name
// - First line is the variable value
// Variables represented as files where filename is name of variable, file first line is a value.
// Returns (Environment, nil) on success or (nil, error) on failure.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	log.Println(env)

	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.Contains(fileName, "=") {
			continue // Skip files with '=' in name as per spec
		}

		// Read the first line of the file
		filePath := filepath.Join(dir, fileName)
		firstLine, isEmpty, err := readFirstLineOnly(filePath)
		if err != nil {
			return nil, err
		}

		env[fileName] = EnvValue{
			Value:      firstLine,
			NeedRemove: isEmpty,
		}

		// Apply to environment
		os.Unsetenv(fileName) // Remove existing if any
		os.Setenv(fileName, firstLine)
	}
	return env, nil
}

// Reads just the first line efficiently.
func readFirstLineOnly(path string) (string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		// File is completely empty
		return "", true, nil
	}

	line := scanner.Text()
	if strings.TrimSpace(line) == "" { // first line contains only whitespace.
		return "", false, nil
	}

	line = strings.ReplaceAll(line, "\x00", "\n")
	return strings.TrimRight(line, " \t"), false, nil
}
