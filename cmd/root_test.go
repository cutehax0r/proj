package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRootCommandWithoutArguments(t *testing.T) {
	originalArgs := os.Args
	originalStdout := os.Stdout

	defer func() {
		os.Args = originalArgs
		os.Stdout = originalStdout
	}()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Args = []string{"proj"}

	os.Stdout = writer

	done := make(chan struct{})
	var output string

	go func() {
		defer close(done)
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, reader)
		if err != nil {
			return
		}
		output = buf.String()
	}()

	rootCmd.SetOutput(writer)
	rootCmd.SetArgs([]string{})
	err = rootCmd.Execute()
	writer.Close()

	<-done

	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "contains Usage",
			pattern: "Usage",
		},
		{
			name:    "contains Available Commands",
			pattern: "Available Commands",
		},
		{
			name:    "contains Flags",
			pattern: "Flags",
		},
		{
			name:    "contains new command",
			pattern: "new",
		},
		{
			name:    "contains add command",
			pattern: "add",
		},
		{
			name:    "contains help flag",
			pattern: "help",
		},
		{
			name:    "contains no-write flag",
			pattern: "no-write",
		},
	}

	for _, tt := range tests {
		if !strings.Contains(output, tt.pattern) {
			t.Errorf("Expected help text to contain %q, but it didn't. Output:\n%s", tt.pattern, output)
		}
	}
}
