# Scripts Test Template

This template demonstrates the script execution capabilities of proj.

## Overview

This template tests script execution at three levels:
1. **Global Level** - Scripts that run before/after all template operations
2. **Template Level** - Scripts that run before/after a specific template
3. **Definition Level** - Scripts that run before/after a definition within a template

## Script Execution Order

Scripts execute in the following order:

1. `global_before.lua` - Sets global variables and initializes state
2. `template_before.lua` - Sets template-level variables
3. `definition_before.lua` - Sets definition-specific variables
4. **File operations execute**
5. `definition_after.lua` - Computes definition-specific paths
6. `template_after.lua` - Finalizes template computations
7. `global_after.lua` - Verifies execution and logs completion

## Script Capabilities

### Variable Management
- Scripts can set variables: `proj.variables.key = "value"`
- Variables set in earlier scripts are available in later scripts
- Variables can be used in file paths and templates

### Logging Functions
- `proj.logDebug(message)` - Debug level logging
- `proj.logInfo(message)` - Info level logging
- `proj.logWarn(message)` - Warning level logging
- `proj.logError(message)` - Error level logging

### Access to Context
- `proj.paths` - Access to template and working paths
- `proj.requirements` - Access to requirement specifications
- `proj.variables.targetName` - The name of the target project
- `proj.noWrite` - Whether running in no-write mode

## Test Verification

The scripts in this template:
1. Track their execution order in a table
2. Set variables at each stage
3. Log their progress
4. Compute paths using variables
5. Verify all stages completed successfully

## Usage

```bash
proj new scripts my-test-project
```

This will create a project with files generated using variables computed by the scripts.
