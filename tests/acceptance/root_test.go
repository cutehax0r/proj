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
		ExpectError("Setup logging").
		ExpectError("Log Level").
		ExpectError("DEBUG")
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
		ExpectError("Configuration loaded").
		ExpectError("no-write").
		ExpectError("true")
}

func TestRootCommand_UsesDefaultConfig(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-l", "3").
		ExpectError("global configuration read").
		ExpectError("testdata/config/proj/proj.yml")
}

func TestRootCommand_UsesConfigFileFlag(t *testing.T) {
	ctx := Setup(t)

	projRoot, _ := filepath.Abs("../..")
	configPath := filepath.Join(projRoot, "testdata/config/test-default.yml")

	ctx.Run("-g", configPath, "-l", "3").
		ExpectError("global configuration read").
		ExpectError(configPath)
}

func TestRootCommand_ErrorsOnMissingConfigFile(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("-g", "/nonexistent/config.yml", "-l", "3").
		ExpectError("Global configuration load failure").
		ExpectError("no such file or directory")
}

func TestRootCommand_ErrorsOnInvalidConfigFile(t *testing.T) {
	ctx := Setup(t)

	projRoot, _ := filepath.Abs("../..")
	invalidConfigPath := filepath.Join(projRoot, "testdata/config/invalid.yml")

	ctx.Run("-g", invalidConfigPath, "-l", "3").
		ExpectError("Global configuration load failure").
		ExpectError("While parsing config")
}
