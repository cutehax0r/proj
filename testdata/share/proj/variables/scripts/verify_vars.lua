-- Verify that variables from command line are accessible
proj.logInfo("=== VARIABLES TEST SCRIPT ===")
proj.logInfo("Author: " .. (proj.variables.author or "NOT SET"))
proj.logInfo("Project Title: " .. (proj.variables.project_title or "NOT SET"))
proj.logInfo("Description: " .. (proj.variables.description or "NOT SET"))

-- Test that a script can set a missing required variable if needed
-- This allows TestVariables_ScriptCanSetMissingRequiredVariable to work
-- Only set author, not project_title, so missing required variable tests still fail
if not proj.variables.author then
    proj.variables.author = "SetByScript"
    proj.logInfo("Set author via script fallback")
end
