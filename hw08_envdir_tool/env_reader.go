package main

import (
	"bytes"
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
	// Initialize Viper (but skip auto-parsing) with //"github.com/spf13/viper"
	// v := viper.New()
	// v.SetConfigType("env")

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
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", fileName, err)
			return nil, err
		}

		value := strings.TrimSpace(readFirstLine(content))
		needRemove := false
		if len(value) == 0 {
			needRemove = true // Empty file case
		} else {
			value = strings.TrimRight(value, " \t\r") // Normal case - take first line
			if value == "" {
				needRemove = true
			}
		}

		env[fileName] = EnvValue{
			Value:      value,
			NeedRemove: needRemove,
		}

		// Set the VARIABLE (filename) -> VALUE in Viper
		// v.Set(file.Name(), value) // Filename = Key, First line = Value

		// Apply to environment
		os.Unsetenv(fileName) // Remove existing if any
		os.Setenv(fileName, value)
	}

	return env, nil
}

func readFirstLine(content []byte) string {
	// Replace null bytes with newlines as per spec
	content = bytes.ReplaceAll(content, []byte{0x00}, []byte{'\n'})

	// Find first newline
	end := bytes.IndexByte(content, '\n')
	if end == -1 {
		end = len(content)
	}

	content = content[:end]
	normalized := string(content)
	// Normalize all line endings to `\n` first
	// normalized := strings.ReplaceAll(string(content), "\r\n", "\n") // Windows
	// normalized = strings.ReplaceAll(normalized, "\r", "\n")         // Old Mac
	// normalized = strings.TrimRight(normalized, " \t")               // following Task Spec

	// fmt.Println(" ----------", normalized, "--------- ")
	// Now split and take the first line
	return strings.Split(normalized, "\n")[0]
}
