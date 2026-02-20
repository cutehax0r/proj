package acceptance

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAddCommand_VariableResolution_CLIOverride(t *testing.T) {
	ctx, projectDir := SetupProjectFromTemplate(t, "varproject")

	// Global config has: version=1.0.0, environment=production, appname=DefaultApp
	// Project .proj/proj.yml has: version=1.0.0, environment=staging, appname=TestApp
	// CLI sets: version=2.0.0
	ctx.Run("add", "config", "cli-override", "-v", "version=2.0.0").
		ExpectExitCode(0)

	configPath := filepath.Join(projectDir, "config-cli-override.txt")
	VerifyFileExists(t, configPath)

	content := ReadFileString(t, configPath)

	if !strings.Contains(content, "version=2.0.0") {
		t.Errorf("Expected config to contain version=2.0.0 from CLI override, got:\n%s", content)
	}

	if !strings.Contains(content, "environment=staging") {
		t.Errorf("Expected config to contain environment=staging from project .proj/proj.yml, got:\n%s", content)
	}

	if !strings.Contains(content, "appname=TestApp") {
		t.Errorf("Expected config to contain appname=TestApp from project .proj/proj.yml, got:\n%s", content)
	}
}

func TestAddCommand_VariableResolution_WithoutCLIOverride(t *testing.T) {
	// Should use: Project > Template > Global > Default
	ctx, projectDir := SetupProjectFromTemplate(t, "varproject")

	ctx.Run("add", "config", "no-override").
		ExpectExitCode(0)

	configPath := filepath.Join(projectDir, "config-no-override.txt")
	VerifyFileExists(t, configPath)

	content := ReadFileString(t, configPath)

	if !strings.Contains(content, "version=1.0.0") {
		t.Errorf("Expected config to contain version=1.0.0 from project .proj/proj.yml, got:\n%s", content)
	}

	if !strings.Contains(content, "environment=staging") {
		t.Errorf("Expected config to contain environment=staging from project .proj/proj.yml, got:\n%s", content)
	}

	if !strings.Contains(content, "appname=TestApp") {
		t.Errorf("Expected config to contain appname=TestApp from project .proj/proj.yml, got:\n%s", content)
	}
}

func TestAddCommand_VariableResolution_ProjectOverridesTemplate(t *testing.T) {
	// This test uses SetupWithConfig to ensure we have the variabletest global config
	ctx := SetupWithConfig(t, "variabletest")

	srcDir := filepath.Join(ProjRoot(), "testdata", "projects", "varproject")
	projectDir := filepath.Join(ctx.TempDir, "varproject")
	copyDir(t, srcDir, projectDir)

	ctx.TempDir = projectDir

	// Global config has: version=1.0.0, environment=production
	// Project has: version=1.0.0, environment=staging (from when we created it)
	ctx.Run("add", "config", "test").
		ExpectExitCode(0)

	configPath := filepath.Join(projectDir, "config-TestApp.txt")
	VerifyFileExists(t, configPath)

	content := ReadFileString(t, configPath)

	if !strings.Contains(content, "environment=staging") {
		t.Errorf("Expected config to contain environment=staging from project, got:\n%s", content)
	}

	if !strings.Contains(content, "appname=TestApp") {
		t.Errorf("Expected config to contain appname=TestApp from project, got:\n%s", content)
	}

	if strings.Contains(content, "appname=DefaultApp") {
		t.Errorf("Expected NOT to find appname=DefaultApp from global config, got:\n%s", content)
	}
}
