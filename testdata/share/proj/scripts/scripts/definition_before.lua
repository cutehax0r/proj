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
local target = proj.variables.targetname or "unknown"
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

-- Verify proj module structure
proj.logInfo("=== PROJ MODULE STRUCTURE CHECK ===")
proj.logInfo("PROJ_NOWRITE: " .. tostring(proj.noWrite))
proj.logInfo("PROJ_ISLOCAL: " .. tostring(proj.requirements.isLocal))

-- Check paths table
if proj.paths then
    proj.logInfo("PROJ_PATHS_TEMPLATEPATH: " .. tostring(proj.paths.templatePath))
    proj.logInfo("PROJ_PATHS_TARGETPATH: " .. tostring(proj.paths.targetPath))
    proj.logInfo("PROJ_PATHS_TEMPLATECONFIGPATH: " .. tostring(proj.paths.templateConfigPath))
else
    proj.logError("PROJ_PATHS: nil")
end

-- Check files table
if proj.files then
    local count = 0
    for _ in pairs(proj.files) do count = count + 1 end
    proj.logInfo("PROJ_FILES_COUNT: " .. count)
else
    proj.logInfo("PROJ_FILES_COUNT: 0")
end

-- Check variables accessibility
if proj.variables then
    proj.logInfo("PROJ_VARIABLES_ACCESSIBLE: true")
else
    proj.logInfo("PROJ_VARIABLES_ACCESSIBLE: false")
end

-- Check logging functions
if proj.logDebug and proj.logInfo and proj.logWarn and proj.logError then
    proj.logInfo("PROJ_LOG_FUNCTIONS: available")
else
    proj.logInfo("PROJ_LOG_FUNCTIONS: missing")
end

proj.logInfo("=== DEFINITION BEFORE SCRIPT END ===")
