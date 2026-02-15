//go:build acceptance

package acceptance

import (
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
