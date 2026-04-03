package acceptance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("update", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj update [target] [flags]").
		ExpectOutput("Flags:").
		ExpectOutput("--template-root").
		ExpectOutput("--force").
		ExpectOutput("Global Flags:")
}

func TestUpdateCommand_NonGitRepo(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templatePath, "proj.yml"), []byte("name: test\n"), 0644); err != nil {
		t.Fatalf("Failed to create proj.yml: %v", err)
	}

	ctx.Run("update", "-s", ctx.TempDir).
		ExpectExitCode(0).
		ExpectError("not from git")
}

func TestUpdateCommand_CleanGitNoUpdates(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "test-template").
		ExpectExitCode(0)

	ctx.Run("update", "-s", ctx.TempDir).
		ExpectExitCode(0).
		ExpectError("updated")
}

func TestUpdateCommand_DirtyGitNoUpdates(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "test-template").
		ExpectExitCode(0)

	templatePath := filepath.Join(ctx.TempDir, "test-template")
	if err := os.WriteFile(filepath.Join(templatePath, "modified.txt"), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to create modified file: %v", err)
	}

	ctx.Run("update", "-s", ctx.TempDir).
		ExpectExitCode(0).
		ExpectError("local changes")
}

func TestUpdateCommand_DirtyGitWithUpdates(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "test-template").
		ExpectExitCode(0)

	templatePath := filepath.Join(ctx.TempDir, "test-template")
	if err := os.WriteFile(filepath.Join(templatePath, "modified.txt"), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to create modified file: %v", err)
	}

	ctx.Run("update", "-s", ctx.TempDir).
		ExpectExitCode(0).
		ExpectError("local changes")
}

func TestUpdateCommand_DirtyGitWithUpdatesAndForce(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "test-template").
		ExpectExitCode(0)

	templatePath := filepath.Join(ctx.TempDir, "test-template")
	if err := os.WriteFile(filepath.Join(templatePath, "modified.txt"), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to create modified file: %v", err)
	}

	ctx.Run("update", "-s", ctx.TempDir, "--force").
		ExpectExitCode(0).
		ExpectError("updated")

	if _, err := os.Stat(filepath.Join(templatePath, "modified.txt")); !os.IsNotExist(err) {
		t.Errorf("Expected modified file to be removed after force update")
	}
}

func TestUpdateCommand_SpecificTarget(t *testing.T) {
	ctx, gitServer := SetupWithGitServer(t)

	ctx.Run("install", "-s", ctx.TempDir, gitServer.URL, "test-template").
		ExpectExitCode(0)

	ctx.Run("update", "-s", ctx.TempDir, "test-template").
		ExpectExitCode(0).
		ExpectError("updated")
}

func TestUpdateCommand_TargetNotFound(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("update", "-s", ctx.TempDir, "nonexistent").
		ExpectError("not found").
		ExpectExitCode(0)
}
