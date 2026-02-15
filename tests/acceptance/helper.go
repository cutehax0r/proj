package acceptance

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type TestContext struct {
	T          *testing.T
	BinaryPath string
	TempDir    string
	Stdout     string
	Stderr     string
	ExitCode   int
}

func Setup(t *testing.T) *TestContext {
	t.Helper()

	ctx := &TestContext{T: t}
	ctx.findBinary()
	ctx.setupTempDir()
	ctx.setEnvVars()

	return ctx
}

func (ctx *TestContext) findBinary() {
	ctx.T.Helper()

	projRoot, err := filepath.Abs("../..")
	if err != nil {
		ctx.T.Fatalf("Failed to get project root: %v", err)
	}

	binaryName := "proj"
	if runtime.GOOS == "windows" {
		binaryName = "proj.exe"
	}

	ctx.BinaryPath = filepath.Join(projRoot, "bin", binaryName)

	if _, err := os.Stat(ctx.BinaryPath); os.IsNotExist(err) {
		ctx.T.Fatalf("Binary not found at %s. Run 'make build-local' first.", ctx.BinaryPath)
	}
}

func (ctx *TestContext) setupTempDir() {
	ctx.T.Helper()

	tempDir, err := os.MkdirTemp("", "proj-acceptance-*")
	if err != nil {
		ctx.T.Fatalf("Failed to create temp dir: %v", err)
	}
	ctx.TempDir = tempDir

	ctx.T.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
}

func (ctx *TestContext) setEnvVars() {
	ctx.T.Helper()

	projRoot, err := filepath.Abs("../..")
	if err != nil {
		ctx.T.Fatalf("Failed to get project root: %v", err)
	}

	os.Setenv("XDG_CONFIG_HOME", filepath.Join(projRoot, "testdata", "config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(projRoot, "testdata", "share"))
}

func (ctx *TestContext) Run(args ...string) *TestContext {
	ctx.T.Helper()

	cmd := exec.Command(ctx.BinaryPath, args...)
	cmd.Dir = ctx.TempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	ctx.Stdout = stdout.String()
	ctx.Stderr = stderr.String()

	if exitErr, ok := err.(*exec.ExitError); ok {
		ctx.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		ctx.ExitCode = -1
	} else {
		ctx.ExitCode = 0
	}

	return ctx
}

func (ctx *TestContext) ExpectOutput(contains string) *TestContext {
	ctx.T.Helper()

	if !strings.Contains(ctx.Stdout, contains) {
		ctx.T.Errorf("Expected stdout to contain %q, but got:\n%s", contains, ctx.Stdout)
	}
	return ctx
}

func (ctx *TestContext) ExpectError(contains string) *TestContext {
	ctx.T.Helper()

	if !strings.Contains(ctx.Stderr, contains) {
		ctx.T.Errorf("Expected stderr to contain %q, but got:\n%s", contains, ctx.Stderr)
	}
	return ctx
}

func (ctx *TestContext) ExpectExitCode(code int) *TestContext {
	ctx.T.Helper()

	if ctx.ExitCode != code {
		ctx.T.Errorf("Expected exit code %d, but got %d", code, ctx.ExitCode)
	}
	return ctx
}
