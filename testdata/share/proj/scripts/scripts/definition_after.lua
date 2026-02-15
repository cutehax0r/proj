-- Definition After Script
-- This script runs after file operations for the definition
-- It finalizes definition-specific data

proj.logInfo("=== DEFINITION AFTER SCRIPT START ===")

-- Track execution order (append to string)
if proj.variables.script_execution_order then
    proj.variables.script_execution_order = proj.variables.script_execution_order .. ",definition_after"
else
    proj.variables.script_execution_order = "definition_after"
end

-- Set definition completion variables
proj.variables.definition_stage = "finished"

-- Verify the path was computed earlier
proj.logInfo("Verified computed path: " .. (proj.variables.computed_path or "NOT SET"))

-- Demonstrate accessing requirements
proj.logDebug("Requirements isLocal: " .. tostring(proj.requirements.isLocal))

proj.logInfo("=== DEFINITION AFTER SCRIPT END ===")
