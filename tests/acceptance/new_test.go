package acceptance

import (
	"os"
	"path/filepath"
	"strings"
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

func TestNewCommand_FailsWhenTargetPathExists(t *testing.T) {
	ctx := Setup(t)

	existingDir := filepath.Join(ctx.TempDir, "foo")
	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create existing directory: %v", err)
	}

	ctx.Run("new", "static", "foo").
		ExpectExitCode(1).
		ExpectError("Target path exists")
}

func TestNewCommand_DefinitionNameFlag(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "static", "myproject", "-d", "new_alt").
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "myproject")
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	content := ReadFileString(t, readmePath)
	if !strings.Contains(content, "ALTERNATE README") {
		t.Errorf("Expected readme to contain 'ALTERNATE README' from new_alt definition, got:\n%s", content)
	}
}

func TestNewCommand_InvalidDefinitionName(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "static", "myproject", "-d", "nonexistent").
		ExpectExitCode(1).
		ExpectError("Definition does not exist")
}

func TestNewCommand_FailsWhenDefinitionIsLocal(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "website", "myproject", "-d", "html").
		ExpectExitCode(1).
		ExpectError("cannot use local-only definition to create new project")
}

func TestNewCommand_VariablePropagation(t *testing.T) {
	ctx := SetupWithConfig(t, "vartest")

	ctx.Run("new", "vartest", "myproject", "-v", "required_var=cli_required", "-v", "cli_var=from_cli", "-v", "override_test=cli_wins").
		ExpectExitCode(0)

	projConfigPath := filepath.Join(ctx.TempDir, "myproject", ".proj", "proj.yml")
	VerifyFileExists(t, projConfigPath)

	content := ReadFileString(t, projConfigPath)

	if !strings.Contains(content, "template_name:") {
		t.Errorf("Expected proj.yml to contain template_name, got:\n%s", content)
	}

	if !strings.Contains(content, "global_var:") || !strings.Contains(content, "from_global_config") {
		t.Errorf("Expected proj.yml to contain global_var from global config, got:\n%s", content)
	}
	if !strings.Contains(content, "another_global:") || !strings.Contains(content, "global_123") {
		t.Errorf("Expected proj.yml to contain another_global from global config, got:\n%s", content)
	}

	if !strings.Contains(content, "required_var:") || !strings.Contains(content, "cli_required_modified_by_before") {
		t.Errorf("Expected proj.yml to contain required_var modified by before script, got:\n%s", content)
	}
	if !strings.Contains(content, "cli_var:") || !strings.Contains(content, "from_cli") {
		t.Errorf("Expected proj.yml to contain cli_var from CLI, got:\n%s", content)
	}

	if !strings.Contains(content, "override_test:") || !strings.Contains(content, "cli_wins") {
		t.Errorf("Expected proj.yml to contain override_test with CLI value, got:\n%s", content)
	}

	if !strings.Contains(content, "optional_var:") || !strings.Contains(content, "after_overrides_default") {
		t.Errorf("Expected proj.yml to contain optional_var overridden by after script, got:\n%s", content)
	}

	if !strings.Contains(content, "script_before:") || !strings.Contains(content, "set_in_before") {
		t.Errorf("Expected proj.yml to contain script_before from before script, got:\n%s", content)
	}
	if !strings.Contains(content, "script_after:") || !strings.Contains(content, "set_in_after") {
		t.Errorf("Expected proj.yml to contain script_after from after script, got:\n%s", content)
	}

	if !strings.Contains(content, "targetName:") || !strings.Contains(content, "myproject") {
		t.Errorf("Expected proj.yml to contain targetName, got:\n%s", content)
	}
	if !strings.Contains(content, "templateName:") || !strings.Contains(content, "vartest") {
		t.Errorf("Expected proj.yml to contain templateName, got:\n%s", content)
	}
	if !strings.Contains(content, "definitionName:") || !strings.Contains(content, "new") {
		t.Errorf("Expected proj.yml to contain definitionName, got:\n%s", content)
	}

	if !strings.Contains(content, "scripts:") {
		t.Errorf("Expected proj.yml to contain empty scripts section, got:\n%s", content)
	}
	if !strings.Contains(content, "definitions:") {
		t.Errorf("Expected proj.yml to contain empty definitions section, got:\n%s", content)
	}
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
