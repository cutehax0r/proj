//go:build acceptance

package acceptance

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestScripts_ExecutesAllScriptsInOrder(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	// Verify all scripts executed by checking log output
	expectedOrder := []string{
		"GLOBAL BEFORE SCRIPT START",
		"GLOBAL BEFORE SCRIPT END",
		"TEMPLATE BEFORE SCRIPT START",
		"TEMPLATE BEFORE SCRIPT END",
		"DEFINITION BEFORE SCRIPT START",
		"DEFINITION BEFORE SCRIPT END",
		"DEFINITION AFTER SCRIPT START",
		"DEFINITION AFTER SCRIPT END",
		"TEMPLATE AFTER SCRIPT START",
		"TEMPLATE AFTER SCRIPT END",
		"GLOBAL AFTER SCRIPT START",
		"GLOBAL AFTER SCRIPT END",
	}

	combinedOutput := ctx.Stdout + ctx.Stderr
	for _, expected := range expectedOrder {
		if !strings.Contains(combinedOutput, expected) {
			t.Errorf("Expected output to contain %q, but it didn't.\nCombined Output:\n%s", expected, combinedOutput)
		}
	}

	// Verify script execution order tracking (comma-separated string)
	expectedOrderStr := "global_before,template_before,definition_before,definition_after,template_after,global_after"
	if !strings.Contains(combinedOutput, expectedOrderStr) {
		t.Errorf("Expected execution order %q, but it wasn't found.\nCombined Output:\n%s", expectedOrderStr, combinedOutput)
	}
}

func TestScripts_SetsVariablesAcrossScripts(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	combinedOutput := ctx.Stdout + ctx.Stderr

	// Verify global variables
	if !strings.Contains(combinedOutput, "Set global_stage to: initialized") {
		t.Error("Expected global_stage to be set to 'initialized'")
	}
	if !strings.Contains(combinedOutput, "global_stage: initialized") {
		t.Error("Expected global_stage to be accessible in global_after")
	}

	// Verify template variables
	if !strings.Contains(combinedOutput, "template_stage: completed") {
		t.Error("Expected template_stage to be 'completed' after template_after runs")
	}

	// Verify definition variables
	if !strings.Contains(combinedOutput, "definition_stage: finished") {
		t.Error("Expected definition_stage to be 'finished'")
	}

	// Verify computed values
	if !strings.Contains(combinedOutput, "computed_path: output/scripttest/results") {
		t.Errorf("Expected computed_path to be 'output/scripttest/results', got:\n%s", combinedOutput)
	}

	// Verify full message construction
	if !strings.Contains(combinedOutput, "Computed full_message: tpl_definition processing - template finished") {
		t.Errorf("Expected full_message to be constructed from template and definition variables, got:\n%s", combinedOutput)
	}
}

func TestScripts_CreatesFilesWithComputedPaths(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	_ = ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "scripttest")

	// Verify file was created at computed path
	computedPath := filepath.Join(projectDir, "output", "scripttest", "results", "summary.txt")
	VerifyFileExists(t, computedPath)

	// Verify content of templated file
	content := ReadFileString(t, computedPath)
	if !strings.Contains(content, "Project: scripttest") {
		t.Errorf("Expected file to contain 'Project: scripttest', got:\n%s", content)
	}
	if !strings.Contains(content, "Computed Path: output/scripttest/results") {
		t.Errorf("Expected file to contain computed path, got:\n%s", content)
	}

	// Verify message file
	messagePath := filepath.Join(projectDir, "message.txt")
	VerifyFileExists(t, messagePath)
	messageContent := ReadFileString(t, messagePath)
	if !strings.Contains(messageContent, "tpl_definition processing") {
		t.Errorf("Expected message.txt to contain full_message, got:\n%s", messageContent)
	}
}

func TestScripts_VerifiesLoggingFunctions(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	combinedOutput := ctx.Stdout + ctx.Stderr

	// Verify log levels are working
	if !strings.Contains(combinedOutput, "=== GLOBAL BEFORE SCRIPT") {
		t.Error("Expected Info level logging from global scripts")
	}

	if !strings.Contains(combinedOutput, "This is a debug message") {
		t.Log("Warning: Debug messages may not appear depending on log level")
	}

	if !strings.Contains(combinedOutput, "This is a warning message") {
		t.Error("Expected Warn level logging")
	}

	// Verify template_after logged at different levels
	if !strings.Contains(combinedOutput, "Computed full_message:") {
		t.Error("Expected Info level logging from template_after")
	}
}

func TestScripts_PreservesVariableState(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	combinedOutput := ctx.Stdout + ctx.Stderr

	// Verify global timestamp is set and preserved
	if !strings.Contains(combinedOutput, "Global timestamp is set:") {
		t.Error("Expected global_timestamp to be accessible in definition_before")
	}

	// Verify template variables are accessible in definition scripts
	if !strings.Contains(combinedOutput, "Accessing global_stage: initialized") {
		t.Error("Expected template_before to access global_stage")
	}

	// Verify definition variables are accessible in template_after
	if !strings.Contains(combinedOutput, "Computed full_message: tpl_") {
		t.Error("Expected template_after to access definition variables")
	}
}

func TestScripts_CopiesNonTemplateFiles(t *testing.T) {
	ctx := SetupWithConfig(t, "scripts")
	_ = ctx.Run("-l", "3", "new", "scripts", "scripttest").ExpectExitCode(0)

	projectDir := filepath.Join(ctx.TempDir, "scripttest")

	// Verify readme.md was copied without parsing
	readmePath := filepath.Join(projectDir, "readme.md")
	VerifyFileExists(t, readmePath)

	// Read the template readme to compare
	templateReadmePath := filepath.Join(ProjRoot(), "testdata/share/proj/scripts/readme.md")
	VerifyFileHash(t, projectDir, "readme.md", templateReadmePath)
}
