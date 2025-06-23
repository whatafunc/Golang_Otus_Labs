package main

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRunCmd(t *testing.T) {
	os.Setenv("ADDED", "from original env")
	defer os.Unsetenv("ADDED")

	env := Environment{
		"FOO": EnvValue{Value: "bar"},
	}

	cmd := []string{"/bin/bash", "-c", "echo $ADDED $FOO"}
	exitCode := RunCmd(cmd, env)

	assert.Equal(t, 0, exitCode)
}
