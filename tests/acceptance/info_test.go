package acceptance

import "testing"

func TestInfoCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("info", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj info [template] [definition] [flags]").
		ExpectOutput("--template-root").
		ExpectOutput("--all").
		ExpectOutput("Global Flags:")
}

func TestInfoCommand_NoArgs_ListsTemplates(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("info").
		ExpectExitCode(0).
		ExpectOutput("static").
		ExpectOutput("website")
}

func TestInfoCommand_Template_ShowsNewAndAddSections(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("info", "static").
		ExpectExitCode(0).
		ExpectOutput("# Template: static").
		ExpectOutput("## New").
		ExpectOutput("new (template)").
		ExpectOutput("new_alt (template)").
		ExpectOutput("## Add").
		ExpectOutput("text (template)")
}

func TestInfoCommand_TemplateDefinition_ShowsTargetsAndVariables(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("info", "website", "new").
		ExpectExitCode(0).
		ExpectOutput("# Template: website").
		ExpectOutput("# Definition: new (template)").
		ExpectOutput("## Files").
		ExpectOutput("src/index.html").
		ExpectOutput("## Variables").
		ExpectOutput("sitename")
}

func TestInfoCommand_NoArgs_InProject_ShowsLocalAndCurrentTemplate(t *testing.T) {
	ctx, _ := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("info").
		ExpectExitCode(0).
		ExpectOutput("# Templates").
		ExpectOutput("website (template)").
		ExpectOutput("# Definitions").
		ExpectOutput("form (local)").
		ExpectOutput("js (local)").
		ExpectOutput("stylesheet (local)")
}

func TestInfoCommand_TemplateDefinition_UsesLocalOverrideInProject(t *testing.T) {
	ctx, _ := SetupProjectFromTemplate(t, "testsite")

	ctx.Run("info", "website", "form").
		ExpectExitCode(0).
		ExpectOutput("# Definition: form (local)").
		ExpectOutput("src/{{.targetname}}-form.html").
		ExpectOutput("formtarget")
}

func TestInfoCommand_AllFlag_ForcesGlobalViewInProject(t *testing.T) {
	ctx, _ := SetupProjectFromTemplate(t, "testsite")

	// When inside a project, 'proj info' shows project context
	// But with --all, it should show global template list
	ctx.Run("info", "--all").
		ExpectExitCode(0).
		ExpectOutput("static").
		ExpectOutput("website")
}

func TestInfoCommand_AllFlagShort_ForcesGlobalViewInProject(t *testing.T) {
	ctx, _ := SetupProjectFromTemplate(t, "testsite")

	// Test short form -a
	ctx.Run("info", "-a").
		ExpectExitCode(0).
		ExpectOutput("static").
		ExpectOutput("website")
}

func TestInfoCommand_TargetPath_FromOutsideProject(t *testing.T) {
	ctx := Setup(t)
	_, projectDir := SetupProjectFromTemplate(t, "testsite")

	// Run info from outside the project, but point to it with -p
	ctx.TempDir = "" // Reset to empty so we're not in the project
	ctx.Run("info", "-p", projectDir).
		ExpectExitCode(0).
		ExpectOutput("# Templates").
		ExpectOutput("website (template)").
		ExpectOutput("form (local)")
}

func TestInfoCommand_TargetPath_AllFlag_FromOutsideProject(t *testing.T) {
	ctx := Setup(t)
	_, projectDir := SetupProjectFromTemplate(t, "testsite")

	// When using -p with a project AND --all, should show global templates
	ctx.TempDir = ""
	ctx.Run("info", "-p", projectDir, "--all").
		ExpectExitCode(0).
		ExpectOutput("static").
		ExpectOutput("website")
}

func TestInfoCommand_TargetPath_NonExistentPath(t *testing.T) {
	ctx := Setup(t)

	// Non-existent path with -p should gracefully show global templates
	// (since no project found at that path, it falls back to global)
	ctx.Run("info", "-p", "/nonexistent/path").
		ExpectExitCode(0).
		ExpectOutput("static").
		ExpectOutput("website")
}

func TestInfoCommand_TemplateRoot_UsesCustomTemplateRoot(t *testing.T) {
	ctx := Setup(t)

	// Use the testdata share directory as template root
	templateRoot := ctx.TempDir

	ctx.Run("info", "-s", templateRoot).
		ExpectExitCode(0)
	// Should show templates from the custom root (which is empty in this case)
}
