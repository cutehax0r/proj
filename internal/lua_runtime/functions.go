package lua_runtime

import (
	"github.com/yuin/gopher-lua"
	"log/slog"
)

var LuaRuntimeFunctions = map[string]lua.LGFunction{}

func init() {
	levels := map[string]func(string, ...any){
		"Debug": slog.Debug,
		"Info":  slog.Info,
		"Warn":  slog.Warn,
		"Error": slog.Error,
	}

	for level, logFn := range levels {
		LuaRuntimeFunctions["log"+level] = createLogFunction(logFn)
	}
}

func createLogFunction(logFn func(string, ...any)) lua.LGFunction {
	// TODO: needs enhancing to support concat style log
	return func(l *lua.LState) int {
		message := l.CheckString(1)
		logFn(message)
		return 0
	}
}

// TODO:
// - dump - print table
// - json - json dump a table
// - yml - yml dump a table
// - text mangling functions: camecase/snakecase/upcase/downcase
// - file path mangling functions (base name, extention, etc)
// - run go 'template parsing' sending along `variables` to be templated in.
