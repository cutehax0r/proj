package acceptance

import (
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
