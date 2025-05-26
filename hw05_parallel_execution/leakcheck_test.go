package hw05parallelexecution

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	fmt.Println("DISABLE_GOLEAK =", os.Getenv("DISABLE_GOLEAK"))

	if os.Getenv("DISABLE_GOLEAK") == "" {
		goleak.VerifyTestMain(m)
	} else {
		os.Exit(m.Run())
	}
}
