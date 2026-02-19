-- Set variables in before script
proj.variables.script_before = "set_in_before"
-- Modify the required_var that came from CLI
proj.variables.required_var = proj.variables.required_var .. "_modified_by_before"
