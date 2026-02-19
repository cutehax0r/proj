package acceptance

import (
	"os"
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
	if !strings.Contains(content, "alert('Hello from notifications!')") {
		t.Errorf("Expected notifications.js to contain alert() with message from project definition, got:\n%s", content)
	}
	if !strings.Contains(content, "Display alert message") {
		t.Errorf("Expected notifications.js to contain 'Display alert message' comment from project override, got:\n%s", content)
	}
}

func TestAddCommand_FallbackToTemplateFiles(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "stylesheet", "print").
		ExpectExitCode(0)

	cssPath := filepath.Join(projectDir, "src", "css", "print.css")
	VerifyFileExists(t, cssPath)

	content := ReadFileString(t, cssPath)
	// This should use the template's src/style.css since .proj/src/style.css doesn't exist
	// Verify it's using the template version (should contain typical CSS content)
	if !strings.Contains(content, "css") && !strings.Contains(content, "color") && !strings.Contains(content, "font") {
		t.Errorf("Expected print.css to contain template CSS content, got:\n%s", content)
	}
}

func TestAddCommand_ProjectFilesTakePrecedence(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("add", "js", "custom", "-v", "alertMessage=Test Alert").
		ExpectExitCode(0)

	jsPath := filepath.Join(projectDir, "src", "js", "custom.js")
	VerifyFileExists(t, jsPath)

	content := ReadFileString(t, jsPath)
	if !strings.Contains(content, "alert('Test Alert')") {
		t.Errorf("Expected custom.js to use project definition with alert(), got:\n%s", content)
	}
}

func TestAddCommand_FailsWhenFileWouldBeOverwritten(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	indexPath := filepath.Join(projectDir, "src", "index.html")
	originalHash := FileSHA1(indexPath)

	ctx.Run("add", "html", "index", "-v", "title=Home").
		ExpectExitCode(1).
		ExpectError("already exists")

	finalHash := FileSHA1(indexPath)
	if originalHash != finalHash {
		t.Errorf("Expected index.html to remain unchanged, but its content was modified")
	}
}

func TestAddCommand_WorksFromNestedDirectory(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "testsite")

	// Create a nested directory within the project
	nestedDir := filepath.Join(projectDir, "src", "foo", "bar")
	err := os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}
	originalDir := ctx.TempDir
	ctx.TempDir = nestedDir

	ctx.Run("add", "html", "about", "-v", "title=About Us").
		ExpectExitCode(0)

	ctx.TempDir = originalDir

	// Verify the file was created relative to project root, not nested directory
	// Should be at src/about.html, not src/foo/bar/about.html
	aboutPath := filepath.Join(projectDir, "src", "about.html")
	VerifyFileExists(t, aboutPath)

	content := ReadFileString(t, aboutPath)
	if !strings.Contains(content, "<title>About Us</title>") {
		t.Errorf("Expected about.html to contain title 'About Us', got:\n%s", content)
	}
	if !strings.Contains(content, "<h1>About Us</h1>") {
		t.Errorf("Expected about.html to contain h1 'About Us', got:\n%s", content)
	}
}
