//go:build acceptance

package acceptance

import (
	"testing"
)

func TestRootCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run().
		ExpectOutput("Usage:").
		ExpectOutput("Flags:").
		ExpectOutput("Available Commands:").
		ExpectOutput("new").
		ExpectOutput("add")
}
