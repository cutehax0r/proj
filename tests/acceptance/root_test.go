//go:build acceptance

// NOTE: These tests use short flag forms (-h, -l, -g, -w) for brevity.
// Long forms (--help, --log-level, --global-config-file, --no-write) are also available.

package acceptance

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run().
		ExpectOutput("Usage:").
		ExpectOutput("Flags:").
		ExpectOutput("Available Commands:").
		ExpectOutput("new").
		ExpectOutput("add")
}

func TestRootCommand_HelpFlag(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-h").
		ExpectOutput("Usage:").
		ExpectOutput("Flags:").
		ExpectOutput("Available Commands:")
}

func TestRootCommand_LogLevelDebug(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-l", "3").
		ExpectOutput("Setup logging").
		ExpectOutput("Log Level").
		ExpectOutput("DEBUG")
}

func TestRootCommand_LogLevelDefault(t *testing.T) {
	ctx := Setup(t)

	ctx.Run().
		ExpectOutput("Usage:")

	if strings.Contains(ctx.Stdout, "Setup logging") {
		t.Error("Expected no debug output with default log level, but found 'Setup logging'")
	}
}

func TestRootCommand_NoWriteFlag(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "-w", "-l", "3", "foo", "bar").
		ExpectOutput("Configuration loaded").
		ExpectOutput("no-write").
		ExpectOutput("true")
}

func TestRootCommand_UsesDefaultConfig(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-l", "3").
		ExpectOutput("global configuration read").
		ExpectOutput("testdata/config/proj/proj.yml")
}

func TestRootCommand_UsesConfigFileFlag(t *testing.T) {
	ctx := Setup(t)

	projRoot, _ := filepath.Abs("../..")
	configPath := filepath.Join(projRoot, "testdata/config/test-default.yml")

	ctx.Run("-g", configPath, "-l", "3").
		ExpectOutput("global configuration read").
		ExpectOutput(configPath)
}

func TestRootCommand_ErrorsOnMissingConfigFile(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-g", "/nonexistent/config.yml", "-l", "3").
		ExpectOutput("Global configuration load failure").
		ExpectOutput("no such file or directory")
}

func TestRootCommand_ErrorsOnInvalidConfigFile(t *testing.T) {
	ctx := Setup(t)

	projRoot, _ := filepath.Abs("../..")
	invalidConfigPath := filepath.Join(projRoot, "testdata/config/invalid.yml")

	ctx.Run("-g", invalidConfigPath, "-l", "3").
		ExpectOutput("Global configuration load failure").
		ExpectOutput("While parsing config")
}
