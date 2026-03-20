package acceptance

import "testing"

func TestInfoCommand_ShowsHelp(t *testing.T) {
	ctx := Setup(t)

	ctx.Run("info", "--help").
		ExpectOutput("Usage:").
		ExpectOutput("proj info [template] [definition] [flags]").
		ExpectOutput("--template-root").
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
