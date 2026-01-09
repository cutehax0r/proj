package lua_runtime

import (
	"github.com/yuin/gopher-lua"
	"log/slog"
)

type Runtime struct {
	Path   string
	Config map[string]any
	Error  error
	state  *lua.LState
}

func NewRuntime(config map[string]any, path string) *Runtime {
	// should also accept 'viper config'
	r := Runtime{
		Path:   path,
		Config: config,
	}
	r.setupExecutionEnvironment()
	return &r
}

func (r *Runtime) Run() error { // take a 'reason' to pass to debug?
	slog.Debug("Executing script", "path", r.Path)
	err := r.state.DoFile(r.Path)
	if err != nil {
		slog.Error("Lua Error", "path", r.Path, "error", err)
		r.Error = err
	}
	slog.Debug("Execution finished", "path", r.Path, "success", err == nil)
	return err
}

func (r *Runtime) setupExecutionEnvironment() {
	r.state = lua.NewState()
	r.state.OpenLibs()

	r.state.PreloadModule("proj", func(l *lua.LState) int {
		mod := l.NewTable()

		// TODO: this variable binding stuff kinda sucks
		keys := make([]string, 0, len(r.Config))
		for k := range r.Config {
			keys = append(keys, k)
		}
		slog.Error("KEYS", "keys", keys)
		
		mod.RawSetString("name", lua.LString(r.Config["name"].(string)))
		mod.RawSetString("kind", lua.LString(r.Config["kind"].(string)))
		mod.RawSetString("dry_run", lua.LBool(r.Config["dry_run"].(bool)))
		mod.RawSetString("global_config_file", lua.LString(r.Config["config"].(string)))
		mod.RawSetString("template_root", lua.LString(r.Config["template_root"].(string)))
		mod.RawSetString("template_path", lua.LString(r.Config["template_path"].(string)))
		mod.RawSetString("target_root", lua.LString(r.Config["target_root"].(string)))
		mod.RawSetString("target_path", lua.LString(r.Config["target_path"].(string)))

		// TODO: pull these from config
		mod.RawSetString("variables", lua.LString("this will be a map"))
		mod.RawSetString("requirements", lua.LString("this will be a map"))
		mod.RawSetString("files", lua.LString("this will be a map"))

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
