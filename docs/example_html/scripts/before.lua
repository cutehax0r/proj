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

print("------------------------")
print "BEFORE Config:"
print("------------------------")
print(dump(VARIABLES))
print("------------------------")



-- access a variable set globally in ~/.config/proj/config.yml in the 'variables' section:
print("email: " .. VARIABLES["email"])

-- access a project template level variable. They are merged with global ones
print("baz: " .. VARIABLES["baz"]) -- prints "true" - as a string because all variables are strings (for now)

-- change a variable and have it propagate through the system
VARIALBES["user"] = "User from Lua"
print("changed user")

-- add a new variable - it'll persist to other stages of execution
VARIABLES["very_cute"] = "sure is"
