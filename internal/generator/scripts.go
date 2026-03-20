package generator

import (
	"errors"
	"log/slog"
	"proj/internal/config"
	"proj/internal/luabridge"
	"strings"
)

func runBeforeScriptsAndValidateVars(luaenv *luabridge.Runtime, scripts []string, reqs *config.RequirementSpec) (map[string]any, error) {
	if err := runScripts(luaenv, scripts); err != nil {
		return nil, err
	}

	vars := luaenv.GetVariables()
	slog.Debug("Variables after before-scripts", slog.Any("vars", vars))

	for _, varspec := range reqs.Variables {
		normalizedName := strings.ToLower(varspec.Name)
		if vars[normalizedName] == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any("Name", varspec.Name))
			slog.Info("All variables", slog.Any("vars", vars))
			return nil, errors.New("required variable not set")
		}
	}
	slog.Debug("All the variables are ready so we can do the work")

	return vars, nil
}

func runAfterScripts(luaenv *luabridge.Runtime, scripts []string) error {
	return runScripts(luaenv, scripts)
}

func runScripts(luaenv *luabridge.Runtime, scripts []string) error {
	for _, script := range scripts {
		if err := luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}

	return nil
}
