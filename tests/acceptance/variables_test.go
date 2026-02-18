package acceptance

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestVariables_SingleVariable(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "variables", "testproject", "-v", "author=John", "-v", "project_title=MyProject").
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "testproject")
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	content := ReadFileString(t, readmePath)
	if !strings.Contains(content, "# MyProject") {
		t.Errorf("Expected readme to contain '# MyProject', got:\n%s", content)
	}
	if !strings.Contains(content, "**Author:** John") {
		t.Errorf("Expected readme to contain '**Author:** John', got:\n%s", content)
	}
	if !strings.Contains(content, "**Description:** No description provided") {
		t.Errorf("Expected readme to contain default description, got:\n%s", content)
	}
}

func TestVariables_VariablesWithSpaces(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "variables", "testproject", "-v", "author=John Doe", "-v", "project_title=My Cool Project", "-v", "description=A project with spaces").
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "testproject")
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	content := ReadFileString(t, readmePath)
	if !strings.Contains(content, "# My Cool Project") {
		t.Errorf("Expected readme to contain '# My Cool Project', got:\n%s", content)
	}
	if !strings.Contains(content, "**Author:** John Doe") {
		t.Errorf("Expected readme to contain '**Author:** John Doe', got:\n%s", content)
	}
	if !strings.Contains(content, "**Description:** A project with spaces") {
		t.Errorf("Expected readme to contain description with spaces, got:\n%s", content)
	}
}

func TestVariables_MissingRequiredVariableFails(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "variables", "testproject", "-v", "author=OnlyAuthor").
		ExpectExitCode(1).
		ExpectError("required variable not set")
}

func TestVariables_DefaultValueUsed(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "variables", "testproject", "-v", "author=TestAuthor", "-v", "project_title=TestTitle").
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "testproject")
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	content := ReadFileString(t, readmePath)
	if !strings.Contains(content, "**Description:** No description provided") {
		t.Errorf("Expected readme to contain default description 'No description provided', got:\n%s", content)
	}
}

func TestVariables_ScriptCanSetMissingRequiredVariable(t *testing.T) {
	ctx := Setup(t)

	// Only set project_title, let script provide author
	ctx.Run("new", "variables", "testproject", "-v", "project_title=ScriptTest").
		ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "testproject")
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	content := ReadFileString(t, readmePath)
	if !strings.Contains(content, "# ScriptTest") {
		t.Errorf("Expected readme to contain '# ScriptTest', got:\n%s", content)
	}
	if !strings.Contains(content, "**Author:** SetByScript") {
		t.Errorf("Expected readme to contain author set by script 'SetByScript', got:\n%s", content)
	}

	// Verify script logged the fallback
	combinedOutput := ctx.Stdout + ctx.Stderr
	if !strings.Contains(combinedOutput, "Set author via script fallback") {
		t.Errorf("Expected script to log setting author via fallback, got:\n%s", combinedOutput)
	}
}

func TestVariables_VariablesVisibleInScripts(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("new", "variables", "testproject", "-v", "author=ScriptTest", "-v", "project_title=ScriptTitle", "-v", "description=ScriptDesc").
		ExpectExitCode(0)

	combinedOutput := ctx.Stdout + ctx.Stderr

	// Verify script logged variable values
	if !strings.Contains(combinedOutput, "Author: ScriptTest") {
		t.Errorf("Expected script to log author value, got:\n%s", combinedOutput)
	}
	if !strings.Contains(combinedOutput, "Project Title: ScriptTitle") {
		t.Errorf("Expected script to log project_title value, got:\n%s", combinedOutput)
	}
	if !strings.Contains(combinedOutput, "Description: ScriptDesc") {
		t.Errorf("Expected script to log description value, got:\n%s", combinedOutput)
	}
}
