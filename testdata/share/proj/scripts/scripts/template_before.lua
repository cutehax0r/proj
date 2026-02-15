-- Template Before Script
-- This script runs before the template definition is processed
-- It can access global variables and set template-level variables

proj.logInfo("=== TEMPLATE BEFORE SCRIPT START ===")

-- Track execution order (append to string)
if proj.variables.script_execution_order then
    proj.variables.script_execution_order = proj.variables.script_execution_order .. ",template_before"
else
    proj.variables.script_execution_order = "template_before"
end

-- Set template-level variables
proj.variables.template_stage = "processing"
proj.variables.template_prefix = "tpl_"
proj.variables.template_message = "Hello from template"

-- Access global variables set by global_before
proj.logInfo("Accessing global_stage: " .. (proj.variables.global_stage or "NOT SET"))

-- Demonstrate accessing paths
proj.logDebug("Template path: " .. (proj.paths.templateRoot or "NOT SET"))
proj.logDebug("Target path: " .. (proj.paths.targetPath or "NOT SET"))

proj.logInfo("=== TEMPLATE BEFORE SCRIPT END ===")
