//go:build acceptance

package acceptance

import (
	"path/filepath"
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
