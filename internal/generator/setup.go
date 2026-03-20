package generator

import (
	"errors"
	"log/slog"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
)

func loadAdderState(cfg *Config, p *paths.Paths, reqs *config.RequirementSpec) (map[string]any, config.ScriptSpec, []config.FileSpec, *luabridge.Runtime, error) {
	if !reqs.Local {
		slog.Error("Cannot use non-local definition to add to existing project", slog.String("definition-name", cfg.DefinitionName))
		return nil, config.ScriptSpec{}, nil, nil, errors.New("cannot use non-local definition to add to existing project")
	}

	vars, scripts, files, luaenv, err := loadGeneratorArtifacts(cfg, p, reqs)
	if err != nil {
		return nil, config.ScriptSpec{}, nil, nil, err
	}

	return vars, scripts, files, luaenv, nil
}

func loadCreatorState(cfg *Config, p *paths.Paths, reqs *config.RequirementSpec) (map[string]any, config.ScriptSpec, []config.FileSpec, *luabridge.Runtime, error) {
	if reqs.Local {
		slog.Error("cannot use local-only definition to create new project", slog.String("definition-name", cfg.DefinitionName))
		return nil, config.ScriptSpec{}, nil, nil, errors.New("cannot use local-only definition to create new project")
	}

	vars, scripts, files, luaenv, err := loadGeneratorArtifacts(cfg, p, reqs)
	if err != nil {
		return nil, config.ScriptSpec{}, nil, nil, err
	}

	return vars, scripts, files, luaenv, nil
}

func loadRequirements() (*config.RequirementSpec, error) {
	reqs, err := config.NewRequirements()
	if err != nil {
		slog.Error("Failed to load requirements", slog.Any("error", err))
		return nil, err
	}
	slog.Debug("Final Requirements", slog.Any("reqs", reqs))

	return reqs, nil
}

func loadGeneratorArtifacts(cfg *Config, p *paths.Paths, reqs *config.RequirementSpec) (map[string]any, config.ScriptSpec, []config.FileSpec, *luabridge.Runtime, error) {
	var err error

	vars, err := config.BuildVariables(reqs.Variables)
	if err != nil {
		slog.Error("Failed to build variables", slog.Any("error", err))
		return nil, config.ScriptSpec{}, nil, nil, err
	}
	slog.Debug("Final Variables", slog.Any("vars", vars))

	scripts, err := config.NewScriptSpecWithFS(cfg.Fs, p)
	if err != nil {
		slog.Error("Couldn't build scripts", slog.Any("error", err))
		return nil, config.ScriptSpec{}, nil, nil, err
	}
	slog.Debug("Final scripts", slog.Any("scripts", scripts))

	files, err := config.NewFileSpecsWithFS(cfg.Fs, p)
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		return nil, config.ScriptSpec{}, nil, nil, err
	}
	slog.Debug("Final files", slog.Any("files", files))

	luaenv := luabridge.NewRuntime(vars, p, reqs, &files, cfg.NoWrite)

	return vars, scripts, files, luaenv, nil
}
