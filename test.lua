-- takes a "thing" and dumps it.
function dump(t, l)
  l = l or 1

    if type(t) ~= "table" then return tostring(t) end
    local out = "{\n"
    for k, v in pairs(t) do
      for _ = 0, l, 1 do out  = out .. "  " end
        out = out .. k .. " = " .. dump(v, l+1) .. "\n"
    end
    return out .. "}"
end

print("\n\n********************************************************************************")
print("Test lua")
print("n********************************************************************************")

print("\nTest Proj lib:")
Proj = require("proj")
Proj.logDebug("logDebug Hello from test.lua")
Proj.logInfo("logInfo Hello from test.lua")
Proj.logWarn("logWarn Hello from test.lua")
Proj.logError("logError Hello from test.lua")

print("\nTest importing:")
Imported = require("import")
Imported.hello()

print("\nTest defined function:")
print(dump(Proj))

print("\n********************************************************************************\n\n")
