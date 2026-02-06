package luabridge

import (
	"fmt"
	"log/slog"
	"proj/internal/config"
	"proj/internal/paths"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type Runtime struct {
	Variables    map[string]any
	Paths        *paths.Paths
	Requirements *config.RequirementSpec
	Files        *[]config.FileSpec
	NoWrite      bool
	Error        error
	// need cli parts: template/target/defn
	state *lua.LState
}

func NewRuntime(variables map[string]any, paths *paths.Paths, requirements *config.RequirementSpec, files *[]config.FileSpec, nowrite bool) *Runtime {
	slog.Debug("Lua Bridge setup", "variables", variables)
	r := Runtime{
		Variables:    variables,
		Paths:        paths,
		Requirements: requirements,
		Files:        files,
		NoWrite:      nowrite,
	}
	r.setupExecutionEnvironment()
	return &r
}

func (r *Runtime) Run(script string) error {
	slog.Debug("Executing script", "script", script)
	// we expect a full usable path to the script  paths.Resolve*() can help
	err := r.state.DoFile(script)
	if err != nil {
		slog.Error("Lua Error", "script", script, "error", err)
		r.Error = err
	}
	slog.Debug("Execution finished", "script", script, "success", err == nil)
	// maybe have to capture stdout?
	return err
}

func (r *Runtime) setupExecutionEnvironment() {
	r.state = lua.NewState()
	r.state.OpenLibs()

	r.state.PreloadModule("proj", func(l *lua.LState) int {
		mod := l.NewTable()
		mod.RawSetString("noWrite", lua.LBool(r.NoWrite))

		pathstable := r.Paths.ToMap()
		mod.RawSetString("paths", r.toLuaValue(pathstable))

		// setup files
		var filestable []map[string]any
		if r.Files != nil {
			for _, file := range *r.Files {
				filestable = append(filestable, file.ToMap())
			}
		}
		mod.RawSetString("files", r.toLuaValue(filestable))

		// setup Requirements
		reqtable := l.NewTable()
		reqtable.RawSetString("isLocal", lua.LBool(r.Requirements.Local))

		reqvars := l.NewTable()
		for _, v := range r.Requirements.Variables {
			rte := l.NewTable()
			rte.RawSetString("name", r.toLuaValue(v.Name))
			rte.RawSetString("default", r.toLuaValue(v.Default))
			reqvars.Append(rte)
		}
		reqtable.RawSetString("variables", reqvars)
		mod.RawSetString("requirements", reqtable)

		mod.RawSetString("variables", r.toLuaValue(r.Variables))

		// functions read from ./functions.go
		for name, fn := range LuaRuntimeFunctions {
			mod.RawSetString(name, l.NewFunction(fn))
		}

		l.Push(mod)
		return 1
	})

	if err := r.state.DoString(`
		require('proj')
		proj = package.loaded.proj
	`); err != nil {
		slog.Error("Failed to setup proj module", slog.Any("error", err))
	}
}

func (r *Runtime) CloseState() {
	r.state.Close()
}

func (r *Runtime) toLuaValue(value any) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case int, int8, int16, int32, int64:
		return lua.LNumber(reflect.ValueOf(v).Convert(reflect.TypeOf(int64(0))).Int())
	case float32, float64:
		return lua.LNumber(reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0))).Float())
	}

	// Use reflection recursively for maps and slices - don't forget the keys
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Map:
		tbl := r.state.NewTable()
		for _, key := range val.MapKeys() {
			mapKey := key.Interface()
			mapValue := val.MapIndex(key).Interface()

			var luaKey string
			switch k := mapKey.(type) {
			case string:
				luaKey = k
			default:
				luaKey = fmt.Sprintf("%v", k)
			}

			tbl.RawSetString(luaKey, r.toLuaValue(mapValue))
		}
		return tbl

	case reflect.Slice, reflect.Array:
		tbl := r.state.NewTable()
		for i := 0; i < val.Len(); i++ {
			tbl.Append(r.toLuaValue(val.Index(i).Interface()))
		}
		return tbl
	}

	// just fall through to a default representation
	return lua.LString(fmt.Sprintf("%v", value))
}

func (r *Runtime) fromLuaValue(value lua.LValue) any {
	if value == lua.LNil {
		return nil
	}

	switch v := value.(type) {
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		// Try to preserve integer type if it's a whole number
		num := float64(v)
		if num == float64(int64(num)) {
			return int64(num)
		}
		return num
	case *lua.LTable:
		// Determine if this is an array or a map by checking for consecutive integer keys
		isArray := true
		maxIndex := 0

		v.ForEach(func(key, val lua.LValue) {
			if lnum, ok := key.(lua.LNumber); ok {
				idx := int(lnum)
				if idx > maxIndex {
					maxIndex = idx
				}
			} else {
				isArray = false
			}
		})

		if isArray && maxIndex > 0 {
			// Convert to slice
			result := make([]any, 0, maxIndex)
			for i := 1; i <= maxIndex; i++ {
				val := v.RawGetInt(i)
				result = append(result, r.fromLuaValue(val))
			}
			return result
		} else {
			// Convert to map
			result := make(map[string]any)
			v.ForEach(func(key, val lua.LValue) {
				var keyStr string
				switch k := key.(type) {
				case lua.LString:
					keyStr = string(k)
				case lua.LNumber:
					keyStr = fmt.Sprintf("%v", k)
				default:
					keyStr = fmt.Sprintf("%v", k)
				}
				result[keyStr] = r.fromLuaValue(val)
			})
			return result
		}
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (r *Runtime) GetVariables() map[string]any {
	result := make(map[string]any)

	packageTable := r.state.GetGlobal("package")
	if packageTable == lua.LNil {
		slog.Warn("package table not found in Lua state")
		return result
	}

	pkgTable, ok := packageTable.(*lua.LTable)
	if !ok {
		slog.Warn("package is not a table in Lua state")
		return result
	}

	loadedTable := pkgTable.RawGetString("loaded")
	if loadedTable == lua.LNil {
		slog.Warn("package.loaded not found")
		return result
	}

	loaded, ok := loadedTable.(*lua.LTable)
	if !ok {
		slog.Warn("package.loaded is not a table")
		return result
	}

	mod := loaded.RawGetString("proj")
	if mod == lua.LNil {
		slog.Warn("proj module not found in package.loaded")
		return result
	}

	modTable, ok := mod.(*lua.LTable)
	if !ok {
		slog.Warn("proj is not a table in Lua state")
		return result
	}

	varTable := modTable.RawGetString("variables")
	if varTable == lua.LNil {
		slog.Warn("variables not found in proj module")
		return result
	}

	varTableObj, ok := varTable.(*lua.LTable)
	if !ok {
		slog.Warn("proj.variables is not a table")
		return result
	}

	varTableObj.ForEach(func(key, val lua.LValue) {
		var keyStr string
		switch k := key.(type) {
		case lua.LString:
			keyStr = string(k)
		case lua.LNumber:
			keyStr = fmt.Sprintf("%v", k)
		default:
			keyStr = fmt.Sprintf("%v", k)
		}
		result[keyStr] = r.fromLuaValue(val)
	})

	slog.Debug("Extracted variables from Lua", slog.Any("variables", result))
	return result
}
