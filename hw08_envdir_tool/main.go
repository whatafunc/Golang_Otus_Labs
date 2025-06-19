package main

import (
	"log"
	"os"
)

func main() {
	dir := os.Args[1]        // Get absolute path to directory from CLI
	env, err := ReadDir(dir) // Make ENV vars
	if err != nil {
		log.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	cmd := os.Args[2:]
	resCode := RunCmd(cmd, env) // Process the CLI cmd with ARGS
	os.Exit(resCode)
}
