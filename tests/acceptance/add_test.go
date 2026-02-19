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
