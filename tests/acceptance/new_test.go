package acceptance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj new <kind> <name> [flags]").
		ExpectOutput("Flags:").
		ExpectOutput("--definition-name").
		ExpectOutput("--target-path").
		ExpectOutput("--target-root").
		ExpectOutput("Global Flags:")
}

func TestNewCommand_RequiresArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new").
		ExpectError("Usage:").
		ExpectError("proj new <kind> <name> [flags]").
		ExpectExitCode(1)
}

func TestNewCommand_RequiresTwoArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "foo").
		ExpectError("accepts 2 arg(s), received 1").
		ExpectError("Usage:").
		ExpectExitCode(1)
}

func TestNewCommand_FailsWhenTemplateNotFound(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "foo", "bar").
		ExpectError("Template configuration load failure").
		ExpectExitCode(1)
}

func TestNewCommand_TargetRootFlag(t *testing.T) {
	ctx := Setup(t)

	customTargetRoot := filepath.Join(ctx.TempDir, "custom-output")
	err := os.MkdirAll(customTargetRoot, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom target root: %v", err)
	}

	ctx.Run("new", "static", "myproject", "-r", customTargetRoot).
		ExpectExitCode(0)

	projectDir := filepath.Join(customTargetRoot, "myproject")
	templateDir := filepath.Join(ProjRoot(), "testdata/share/proj/static")

	VerifyFileHash(t, projectDir, "readme.md", filepath.Join(templateDir, "readme.md"))
}

func TestNewCommand_TargetPathFlag(t *testing.T) {
	ctx := Setup(t)

	customTargetPath := filepath.Join(ctx.TempDir, "custom-project-location")

	ctx.Run("new", "static", "ignored_name", "-p", customTargetPath).
		ExpectExitCode(0)

	templateDir := filepath.Join(ProjRoot(), "testdata/share/proj/static")

	VerifyFileHash(t, customTargetPath, "readme.md", filepath.Join(templateDir, "readme.md"))
}

func TestNewCommand_TemplateRootFlag(t *testing.T) {
	ctx := SetupWithoutDataHome(t)

	customTemplateRoot := filepath.Join(ProjRoot(), "testdata/share/proj")

	ctx.Run("new", "static", "myproject", "-s", customTemplateRoot).
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "myproject")
	templateDir := filepath.Join(ProjRoot(), "testdata/share/proj/static")

	VerifyFileHash(t, projectDir, "readme.md", filepath.Join(templateDir, "readme.md"))
}

func TestNewCommand_TemplatePathFlag(t *testing.T) {
	ctx := SetupWithoutDataHome(t)

	customTemplatePath := filepath.Join(ProjRoot(), "testdata/share/proj/static")

	ctx.Run("new", "ignored_kind", "myproject", "-t", customTemplatePath).
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "myproject")

	VerifyFileHash(t, projectDir, "readme.md", filepath.Join(customTemplatePath, "readme.md"))
}
