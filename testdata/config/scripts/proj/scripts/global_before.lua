-- Global Before Script
-- This script runs before any template processing begins
-- It sets up global variables that will be available throughout

proj.logInfo("=== GLOBAL BEFORE SCRIPT START ===")

-- Set global variables (initialize as simple strings for compatibility)
proj.variables.global_stage = "initialized"
proj.variables.global_timestamp = tostring(os.time())
proj.variables.script_execution_order = "global_before"

proj.logDebug("Set global_stage to: " .. proj.variables.global_stage)
proj.logDebug("Global timestamp: " .. proj.variables.global_timestamp)

proj.logInfo("=== GLOBAL BEFORE SCRIPT END ===")
