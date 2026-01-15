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
	Variables map[string]any
	Paths *paths.Paths
	Requirements *config.RequirementSpec
	NoWrite bool
	Error  error
	// will gain reqs and files
	state  *lua.LState
}

// will gain reqs and files
func NewRuntime(variables map[string]any, paths *paths.Paths, requirements *config.RequirementSpec, nowrite bool) *Runtime {
	slog.Debug("Lua Bridge setup", "variables", variables)
	r := Runtime{
		Variables: variables,
		Paths: paths,
		Requirements: requirements,
		NoWrite: nowrite,
	}
	r.setupExecutionEnvironment()
	return &r
}

func (r *Runtime) Run(script string) error {
	slog.Debug("Executing script", "script", script)
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

		// TODO: this variable binding stuff kinda sucks
		keys := make([]string, 0, len(r.Variables))
		for k := range r.Variables {
			keys = append(keys, k)
		}
		slog.Debug("Go variable keys", "keys", keys)
		
		// mod.RawSetString("name", lua.LString(r.Variables["name"].(string)))
		// mod.RawSetString("template", lua.LString(r.Variables["template"].(string)))
		// mod.RawSetString("definition", lua.LString(r.Variables["definition"].(string)))

		// need to bind:
		//* variables
		//* requirements
		//* files

		// functions read from ./functions.go
		for name, fn := range LuaRuntimeFunctions {
			mod.RawSetString(name, l.NewFunction(fn))
		}

		l.Push(mod)
		return 1
	})
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
	case []any:
		tbl := r.state.NewTable()
		for _, i := range v {
			tbl.Append(r.toLuaValue(i))
		}
		return tbl
	case map[string]any:
		tbl := r.state.NewTable()
		for key, val := range(v) {
			tbl.RawSetString(key, r.toLuaValue(val))
		}
		return tbl
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
