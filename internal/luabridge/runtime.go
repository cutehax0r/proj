package luabridge

import (
	"log/slog"
	"proj/internal/paths"

	lua "github.com/yuin/gopher-lua"
)

type Runtime struct {
	Variables map[string]any
	Paths *paths.Paths
	Nowrite bool
	Error  error
	state  *lua.LState
}

func NewRuntime(variables map[string]any, paths *paths.Paths, nowrite bool) *Runtime {
	r := Runtime{
		Variables: variables,
		Paths: paths,
		Nowrite: nowrite,
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

		// TODO: this variable binding stuff kinda sucks
		keys := make([]string, 0, len(r.Variables))
		for k := range r.Variables {
			keys = append(keys, k)
		}
		slog.Error("KEYS", "keys", keys)
		
		mod.RawSetString("name", lua.LString(r.Variables["name"].(string)))
		mod.RawSetString("template", lua.LString(r.Variables["template"].(string)))
		mod.RawSetString("definition", lua.LString(r.Variables["definition"].(string)))
		mod.RawSetString("nowrite", lua.LBool(r.Variables["nowrite"].(bool)))

		// TODO: pull these from config
		mod.RawSetString("variables", lua.LString("this will be a map"))
		mod.RawSetString("requirements", lua.LString("this will be a requirementspec"))
		mod.RawSetString("paths", lua.LString("paths"))
		mod.RawSetString("files", lua.LString("this will be a []filespec"))

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
