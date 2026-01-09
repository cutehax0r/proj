-- you can make functions:
function dump(t)
    if type(t) ~= "table" then return tostring(t) end
    local out = "{"
    local first = true
    for k, v in pairs(t) do
        if not first then out = out .. ", " end
        out = out .. k .. " = " .. dump(v)
        first = false
    end
    return out .. "}"
end

proj = require("proj")
print(dump(proj))
proj.logDebug("Hello, this is a debug message in before.lua")
