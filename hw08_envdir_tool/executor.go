package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Security fix: resolve full executable path
	// G204: Subprocess launched with a potential tainted input or cmd arguments
	binary, err := exec.LookPath(cmd[0])
	if err != nil {
		return 127 // Command not found
	}

	currentEnv := prepareEnv(env)
	// Execute the command with modified environment
	command := exec.Command(binary, cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = currentEnv // Includes our modified vars

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		return 1
	}
	return 0
}

func prepareEnv(customEnv Environment) []string {
	// Start with a clean copy of the current environment
	currentEnv := os.Environ()

	// Convert to map for easier manipulation
	envMap := make(map[string]string)
	for _, env := range currentEnv {
		if keyVal := strings.SplitN(env, "=", 2); len(keyVal) == 2 {
			envMap[keyVal[0]] = keyVal[1]
		}
	}

	// Apply custom environment changes
	for key, envValue := range customEnv {
		if envValue.NeedRemove {
			delete(envMap, key)
		} else {
			envMap[key] = envValue.Value
		}
	}

	// Convert back to slice
	result := make([]string, 0, len(envMap))
	for key, value := range envMap {
		result = append(result, key+"="+value)
	}

	return result
}
