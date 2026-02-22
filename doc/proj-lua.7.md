proj-lua(7) -- Lua scripting bridge
===================================

## NAME

proj-lua - Lua scripting bridge and available functions

## DESCRIPTION

proj embeds a Lua interpreter for advanced scripting capabilities. Scripts
written in Lua can access the full Lua 5.1 standard library plus proj-specific
extensions.

## STANDARD LIBRARY

The following Lua standard libraries are available:

* `base` - Core Lua functions
* `string` - String manipulation
* `table` - Table manipulation
* `math` - Mathematical functions
* `os` - Operating system interface (limited)
* `io` - Input/output (limited)

## PROJ MODULE

The `proj` module provides proj-specific functionality:

```lua
local proj = require("proj")
```

### Functions

* `proj.logInfo(name)` - Log stuff at info level

## SEE ALSO

proj(1), proj-scripts(7), proj-variables(7), proj-template(7)
