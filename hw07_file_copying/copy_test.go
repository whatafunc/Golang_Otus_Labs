package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	// Setup: Create a test source and output dst file
	testFile := "go.sum"
	testOutFile := "a.sum"

	// Get absolute path to test file
	srcPath, err := filepath.Abs(testFile)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	log.Print("ok, file = ", srcPath)
	// Create temp output file path
	dstPath := filepath.Join(t.TempDir(), testOutFile)

	defer os.Remove(dstPath) // Clean up after test
	// Read hardcoded file to get expected content
	originalContent, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	tests := []struct {
		name          string
		offset        int64
		limit         int64
		expectedError bool
	}{
		{"Full file copy", 0, 0, false},
		{"Partial copy with offset", 5, 10, false},
		{"Offset exceeds file size", int64(len(originalContent) + 1), 0, true},
		{"Limit exceeds remaining bytes", 0, int64(len(originalContent)) + 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(srcPath, dstPath, tt.offset, tt.limit)
			if (err != nil) != tt.expectedError {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.expectedError)
			}

			// Verify copied content (for successful cases)
			if !tt.expectedError {
				copiedContent, err := os.ReadFile(dstPath)
				if err != nil {
					t.Errorf("Failed to read destination file: %v", err)
				}

				// Calculate expected content
				expectedSize := tt.limit
				if tt.limit == 0 {
					expectedSize = int64(len(originalContent)) - tt.offset
				} else if tt.offset+tt.limit > int64(len(originalContent)) {
					expectedSize = int64(len(originalContent)) - tt.offset
				}

				if int64(len(copiedContent)) != expectedSize {
					t.Errorf("Copied file size = %d, want %d", len(copiedContent), expectedSize)
				}
			}
		})
	}
}

func TestSpecialFiles(t *testing.T) {
	t.Run("Should reject device files", func(t *testing.T) {
		err := Copy("/dev/urandom", "output.dat", 0, 0)
		if err == nil {
			t.Error("Expected error for device file")
		}
	})
}

func TestSameFileDetection(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0o644)
	defer os.Remove(tmpFile)

	t.Run("Identical paths", func(t *testing.T) {
		err := Copy(tmpFile, tmpFile, 0, 0)
		if err == nil || !strings.Contains(err.Error(), "destination file is same as source file") {
			t.Errorf("Should reject identical paths")
		}
	})

	t.Run("Hard links", func(t *testing.T) {
		linkPath := filepath.Join(t.TempDir(), "link.txt")
		os.Link(tmpFile, linkPath) // Create hard link
		defer os.Remove(linkPath)
		err := Copy(tmpFile, linkPath, 0, 0)
		if err == nil || !strings.Contains(err.Error(), "destination file is symlinked to source file") {
			t.Errorf("Should detect hard links as same file")
		}
	})
}
