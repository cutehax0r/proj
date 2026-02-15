-- Template After Script
-- This script runs after the template definition is processed
-- It can finalize template-level computations (though files have already been created)

proj.logInfo("=== TEMPLATE AFTER SCRIPT START ===")

-- Track execution order (append to string)
if proj.variables.script_execution_order then
    proj.variables.script_execution_order = proj.variables.script_execution_order .. ",template_after"
else
    proj.variables.script_execution_order = "template_after"
end

-- Set template completion variables
proj.variables.template_stage = "completed"

-- Demonstrate log levels
proj.logDebug("This is a debug message from template_after")
proj.logWarn("This is a warning message from template_after")

proj.logInfo("=== TEMPLATE AFTER SCRIPT END ===")
