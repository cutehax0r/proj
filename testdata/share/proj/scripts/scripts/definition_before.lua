-- Definition Before Script
-- This script runs before file operations for the definition
-- It sets definition-specific variables and computes paths for file creation

proj.logInfo("=== DEFINITION BEFORE SCRIPT START ===")

-- Track execution order (append to string)
if proj.variables.script_execution_order then
    proj.variables.script_execution_order = proj.variables.script_execution_order .. ",definition_before"
else
    proj.variables.script_execution_order = "definition_before"
end

-- Set definition-level variables
proj.variables.definition_stage = "executing"
proj.variables.definition_message = "definition processing"

-- Compute a path using variables and target name (this runs before files are created)
local target = proj.variables.targetName or "unknown"
proj.variables.computed_path = "output/" .. target .. "/results"

proj.logInfo("Computed path: " .. proj.variables.computed_path)

-- Compute full_message using variables from all stages
local prefix = proj.variables.template_prefix or ""
local def_msg = proj.variables.definition_message or ""
proj.variables.full_message = prefix .. def_msg .. " - template finished"

proj.logInfo("Computed full_message: " .. proj.variables.full_message)

-- Access variables from previous scripts
proj.logInfo("Global stage: " .. (proj.variables.global_stage or "NOT SET"))
proj.logInfo("Template stage: " .. (proj.variables.template_stage or "NOT SET"))

-- Demonstrate logging (but don't fail - just log)
if proj.variables.global_timestamp == nil then
    proj.logError("Warning: global_timestamp not set!")
else
    proj.logInfo("Global timestamp is set: " .. proj.variables.global_timestamp)
end

proj.logInfo("=== DEFINITION BEFORE SCRIPT END ===")
