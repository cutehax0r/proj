package acceptance

import (
	"bytes"
	"crypto/sha1"
	"fmt"
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
	return SetupWithConfig(t, "default")
}

func SetupWithConfig(t *testing.T, configName string) *TestContext {
	t.Helper()

	ctx := &TestContext{T: t}
	ctx.findBinary()
	ctx.setupTempDir()
	ctx.setEnvVars(configName)

	return ctx
}

func SetupProject(t *testing.T, templateName, projectName string) (*TestContext, string) {
	t.Helper()

	ctx := SetupWithConfig(t, templateName)
	ctx.Run("new", templateName, projectName).ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, projectName)
	return ctx, projectDir
}

func (ctx *TestContext) findBinary() {
	ctx.T.Helper()

	binaryName := "proj"
	if runtime.GOOS == "windows" {
		binaryName = "proj.exe"
	}

	ctx.BinaryPath = filepath.Join(ProjRoot(), "bin", binaryName)

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

func (ctx *TestContext) setEnvVars(configName string) {
	ctx.T.Helper()

	os.Setenv("PROJ_ROOT", ProjRoot())
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(ProjRoot(), "testdata", "config", configName))
	os.Setenv("XDG_DATA_HOME", filepath.Join(ProjRoot(), "testdata", "share"))
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

func ProjRoot() string {
	projRoot, err := filepath.Abs("../..")
	if err != nil {
		panic("Failed to get project root: " + err.Error())
	}
	return projRoot
}

func FileSHA1(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to read file: " + err.Error())
	}
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func VerifyFileHash(t *testing.T, projectDir, relPath, sourcePath string) {
	t.Helper()
	targetPath := filepath.Join(projectDir, relPath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it does not", relPath)
		return
	}

	expectedHash := FileSHA1(sourcePath)
	actualHash := FileSHA1(targetPath)

	if expectedHash != actualHash {
		t.Errorf("File %s content mismatch: expected hash %s, got %s", relPath, expectedHash, actualHash)
	}
}

func VerifyFileHashEquals(t *testing.T, projectDir, relPath string, expectedHash string) {
	t.Helper()
	targetPath := filepath.Join(projectDir, relPath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it does not", relPath)
		return
	}

	actualHash := FileSHA1(targetPath)

	if expectedHash != actualHash {
		t.Errorf("File %s content mismatch: expected hash %s, got %s", relPath, expectedHash, actualHash)
	}
}

func VerifyFileMode(t *testing.T, projectDir, relPath string, expectedMode os.FileMode) {
	t.Helper()
	targetPath := filepath.Join(projectDir, relPath)

	info, err := os.Stat(targetPath)
	if err != nil {
		t.Errorf("Failed to stat file %s: %v", relPath, err)
		return
	}

	actualMode := info.Mode() & os.ModePerm
	if actualMode != expectedMode {
		t.Errorf("File %s mode mismatch: expected %04o, got %04o", relPath, expectedMode, actualMode)
	}
}

func VerifyFileExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s, but it does not", path)
	}
}

func ReadFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(data)
}
