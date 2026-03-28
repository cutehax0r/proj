package acceptance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("install", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj install [source] [target] [flags]").
		ExpectOutput("Flags:").
		ExpectOutput("--template-git").
		ExpectOutput("--template-root").
		ExpectOutput("Global Flags:")
}

func TestInstallCommand_RequiresArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("install").
		ExpectError("Usage:").
		ExpectError("proj install [source] [target] [flags]").
		ExpectExitCode(1)
}

func TestInstallCommand_FromHTTP(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL).
		ExpectExitCode(0)

	templatePath := filepath.Join(ctx.TempDir, "test-template")
	if _, err := os.Stat(templatePath); err != nil {
		t.Errorf("Expected template at %s, got error: %v", templatePath, err)
	}
}

func TestInstallCommand_HTTPWithCustomTarget(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "my-custom-name").
		ExpectExitCode(0)

	templatePath := filepath.Join(ctx.TempDir, "my-custom-name")
	if _, err := os.Stat(templatePath); err != nil {
		t.Errorf("Expected template at %s, got error: %v", templatePath, err)
	}
}

func TestInstallCommand_AlreadyExists(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL).
		ExpectExitCode(0)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL).
		ExpectExitCode(1).
		ExpectError("already exists")
}
