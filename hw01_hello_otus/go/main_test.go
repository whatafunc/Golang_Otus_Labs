package main

import "testing"

// TestReverseGreeting tests the ReverseGreeting function.
func TestReverseGreeting(t *testing.T) {
	result := reverseGreeting("Hello, OTUS!")
	expected := "!SUTO ,olleH"

	if result != expected {
		t.Errorf("Reversal result = %s; want %s", result, expected)
	}
}
