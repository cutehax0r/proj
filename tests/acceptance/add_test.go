package acceptance

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAddCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("add", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj add <kind> <name> [flags]").
		ExpectOutput("Flags:").
		ExpectOutput("--target-path").
		ExpectOutput("--template-path").
		ExpectOutput("Global Flags:")
}

func TestAddCommand_RequiresArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("add").
		ExpectError("Usage:").
		ExpectError("proj add <kind> <name> [flags]").
		ExpectExitCode(1)
}

func TestAddCommand_RequiresTwoArgs(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("add", "foo").
		ExpectError("accepts 2 arg(s), received 1").
		ExpectError("Usage:").
		ExpectExitCode(1)
}

func TestAddCommand_FailsWhenNotInProjDirectory(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("add", "foo", "bar").
		ExpectError("Failed to setup configuration").
		ExpectError("not in a proj directory").
		ExpectExitCode(1)
}

func TestAddCommand_AddHtmlPage(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "html", "contact", "-v", "title=Contact Us").
		ExpectExitCode(0)

	htmlPath := filepath.Join(projectDir, "src", "contact.html")
	VerifyFileExists(t, htmlPath)

	content := ReadFileString(t, htmlPath)
	if !strings.Contains(content, "<title>Contact Us</title>") {
		t.Errorf("Expected contact.html to contain title 'Contact Us', got:\n%s", content)
	}
	if !strings.Contains(content, "<h1>Contact Us</h1>") {
		t.Errorf("Expected contact.html to contain h1 'Contact Us', got:\n%s", content)
	}
}

func TestAddCommand_FailsWhenDefinitionNotLocal(t *testing.T) {
	ctx, _ := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "new", "somepage", "-v", "sitename=Test Page").
		ExpectExitCode(1).
		ExpectError("cannot use non-local definition to add to existing project")
}

func TestAddCommand_AddFormWithProjectDefinition(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "form", "contact", "-v", "formTarget=/api/submit").
		ExpectExitCode(0)

	formPath := filepath.Join(projectDir, "src", "contact-form.html")
	VerifyFileExists(t, formPath)

	content := ReadFileString(t, formPath)
	if !strings.Contains(content, `action="/api/submit"`) {
		t.Errorf("Expected contact-form.html to contain action='/api/submit', got:\n%s", content)
	}
	if !strings.Contains(content, "<h1>contact Form</h1>") {
		t.Errorf("Expected contact-form.html to contain 'contact Form' header, got:\n%s", content)
	}
	if !strings.Contains(content, `method="post"`) {
		t.Errorf("Expected contact-form.html to contain method='post', got:\n%s", content)
	}
}

func TestAddCommand_OverridesTemplateDefinition(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "js", "notifications", "-v", "alertMessage=Hello from notifications!").
		ExpectExitCode(0)

	jsPath := filepath.Join(projectDir, "src", "js", "notifications.js")
	VerifyFileExists(t, jsPath)

	content := ReadFileString(t, jsPath)
	// This should use the project-local definition which has the alert() call
	if !strings.Contains(content, "alert('Hello from notifications!')") {
		t.Errorf("Expected notifications.js to contain alert() with message from project definition, got:\n%s", content)
	}
	// Verify it's using the project override, not the template default
	if !strings.Contains(content, "Display alert message") {
		t.Errorf("Expected notifications.js to contain 'Display alert message' comment from project override, got:\n%s", content)
	}
}
