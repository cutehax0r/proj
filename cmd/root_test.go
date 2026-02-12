package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureCommandOutput(args []string) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		return ""
	}

	originalStdout := os.Stdout
	os.Stdout = writer

	done := make(chan string)
	go func() {
		buf := &bytes.Buffer{}
		io.Copy(buf, reader)
		done <- buf.String()
	}()

	rootCmd.SetOutput(writer)
	rootCmd.SetArgs(args)
	rootCmd.Execute()
	writer.Close()
	os.Stdout = originalStdout

	return <-done
}

func TestRootCommandWithoutArguments(t *testing.T) {
	output := captureCommandOutput([]string{})

	patterns := []string{"Usage", "Available Commands", "Flags", "new", "add", "help", "no-write"}

	for _, pattern := range patterns {
		if !strings.Contains(output, pattern) {
			t.Errorf("Expected help text to contain %q", pattern)
		}
	}
}

func TestNewCommandWithoutArguments(t *testing.T) {
	output := captureCommandOutput([]string{"new"})

	patterns := []string{"Usage", "Flags", "help", "set-variable", "template-root"}

	for _, pattern := range patterns {
		if !strings.Contains(output, pattern) {
			t.Errorf("Expected new command help to contain %q", pattern)
		}
	}
}

func TestAddCommandWithoutArguments(t *testing.T) {
	output := captureCommandOutput([]string{"add"})

	patterns := []string{"Usage", "Flags", "help", "set-variable", "template-root"}

	for _, pattern := range patterns {
		if !strings.Contains(output, pattern) {
			t.Errorf("Expected add command help to contain %q", pattern)
		}
	}
}
