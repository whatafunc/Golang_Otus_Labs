package logger

import (
	"os"
	"testing"
)

func captureOutput(f func(), target **os.File) string {
	reader, writer, _ := os.Pipe()
	orig := *target
	*target = writer
	outC := make(chan string)
	go func() {
		var buf [1024]byte
		n, _ := reader.Read(buf[:])
		outC <- string(buf[:n])
	}()
	f()
	writer.Close()
	*target = orig
	return <-outC
}

func TestLogger_InfoLevel(t *testing.T) {
	log := New("info")
	infoOut := captureOutput(func() {
		log.Info("info message")
	}, &os.Stdout)
	errOut := captureOutput(func() {
		log.Error("error message")
	}, &os.Stderr)

	if infoOut == "" || infoOut != "info message\n" {
		t.Errorf("Info message not found in stdout for info level, got: %q", infoOut)
	}
	if errOut == "" || errOut != "error message\n" {
		t.Errorf("Error message not found in stderr for info level, got: %q", errOut)
	}
}

func TestLogger_ErrorLevel(t *testing.T) {
	log := New("error")
	infoOut := captureOutput(func() {
		log.Info("info message")
	}, &os.Stdout)
	errOut := captureOutput(func() {
		log.Error("error message")
	}, &os.Stderr)

	if infoOut != "" {
		t.Errorf("Expected no info output for error level, got: %q", infoOut)
	}
	if errOut == "" || errOut != "error message\n" {
		t.Errorf("Error message not found in stderr for error level, got: %q", errOut)
	}
}
