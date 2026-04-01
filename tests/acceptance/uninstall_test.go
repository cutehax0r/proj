package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUninstallCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("uninstall", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj uninstall [target] [flags]").
		ExpectOutput("Flags:").
		ExpectOutput("--template-root").
		ExpectOutput("--force").
		ExpectOutput("Global Flags:")
}

func TestUninstallCommand_RequiresArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("uninstall").
		ExpectError("Usage:").
		ExpectError("proj uninstall [target] [flags]").
		ExpectExitCode(1)
}

func TestUninstallCommand_TargetNotFound(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName).
		ExpectError("not a git repository").
		ExpectExitCode(1)
}

func TestUninstallCommand_NotAGitRepo(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName).
		ExpectError("not a git repository").
		ExpectExitCode(1)
}

func TestUninstallCommand_NotAGitRepoWithForce(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName, "--force").
		ExpectExitCode(0)

	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Errorf("Expected template to be removed, but it still exists at %s", templatePath)
	}
}

func TestUninstallCommand_DirtyRepo(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	srcDir := filepath.Join(ProjRoot(), "testdata", "projects", "testsite")
	copyDir(t, srcDir, templatePath)

	if err := runGitCmd(templatePath, "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.name", "Test"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "add", "."); err != nil {
		t.Fatalf("Failed to add files: %v", err)
	}
	if err := runGitCmd(templatePath, "commit", "-m", "initial"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	if err := os.WriteFile(filepath.Join(templatePath, "newfile.txt"), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName).
		ExpectError("uncommitted local changes").
		ExpectExitCode(1)
}

func TestUninstallCommand_DirtyRepoWithForce(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	srcDir := filepath.Join(ProjRoot(), "testdata", "projects", "testsite")
	copyDir(t, srcDir, templatePath)

	if err := runGitCmd(templatePath, "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.name", "Test"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "add", "."); err != nil {
		t.Fatalf("Failed to add files: %v", err)
	}
	if err := runGitCmd(templatePath, "commit", "-m", "initial"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	if err := os.WriteFile(filepath.Join(templatePath, "newfile.txt"), []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName, "--force").
		ExpectExitCode(0)

	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Errorf("Expected template to be removed, but it still exists at %s", templatePath)
	}
}

func TestUninstallCommand_CleanRepo(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	srcDir := filepath.Join(ProjRoot(), "testdata", "projects", "testsite")
	copyDir(t, srcDir, templatePath)

	if err := runGitCmd(templatePath, "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.name", "Test"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "add", "."); err != nil {
		t.Fatalf("Failed to add files: %v", err)
	}
	if err := runGitCmd(templatePath, "commit", "-m", "initial"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName).
		ExpectExitCode(0)

	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Errorf("Expected template to be removed, but it still exists at %s", templatePath)
	}
}

func TestUninstallCommand_NoWriteDryRun(t *testing.T) {
	ctx := Setup(t)

	templateName := "test-template"
	templatePath := filepath.Join(ctx.TempDir, templateName)
	srcDir := filepath.Join(ProjRoot(), "testdata", "projects", "testsite")
	copyDir(t, srcDir, templatePath)

	if err := runGitCmd(templatePath, "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "config", "user.name", "Test"); err != nil {
		t.Fatalf("Failed to config git: %v", err)
	}
	if err := runGitCmd(templatePath, "add", "."); err != nil {
		t.Fatalf("Failed to add files: %v", err)
	}
	if err := runGitCmd(templatePath, "commit", "-m", "initial"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	ctx.Run("uninstall", "-s", ctx.TempDir, templateName, "--no-write").
		ExpectError("Dry run").
		ExpectExitCode(0)

	if _, err := os.Stat(templatePath); err != nil {
		t.Errorf("Expected template to still exist (dry run), but it was removed at %s", templatePath)
	}
}

func runGitCmd(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}
