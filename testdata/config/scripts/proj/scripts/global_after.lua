-- Global After Script
-- This script runs after all template processing is complete
-- It verifies that all scripts ran in the correct order

proj.logInfo("=== GLOBAL AFTER SCRIPT START ===")

-- Track that this script ran (append to string)
if proj.variables.script_execution_order then
    proj.variables.script_execution_order = proj.variables.script_execution_order .. ",global_after"
else
    proj.variables.script_execution_order = "global_after"
end

-- Log execution order
proj.logInfo("Script execution order: " .. (proj.variables.script_execution_order or "NOT TRACKED"))

-- Verify variables set by other scripts
proj.logInfo("Final variable values:")
proj.logInfo("  global_stage: " .. (proj.variables.global_stage or "NOT SET"))
proj.logInfo("  template_stage: " .. (proj.variables.template_stage or "NOT SET"))
proj.logInfo("  definition_stage: " .. (proj.variables.definition_stage or "NOT SET"))
proj.logInfo("  computed_path: " .. (proj.variables.computed_path or "NOT SET"))
proj.logInfo("  full_message: " .. (proj.variables.full_message or "NOT SET"))

proj.logInfo("=== GLOBAL AFTER SCRIPT END ===")
