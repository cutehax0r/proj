package creator

import (
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
)

type Creator struct {
	paths *paths.Paths
	reqs *config.RequirementSpec
	scripts config.ScriptSpec
	files []config.FileSpec
	vars map[string]any
	luaenv *luabridge.Runtime
	noWrite bool
}

func NewCreator(args []string) (*Creator, error) {
	gen := &Creator{}
	return gen, nil
}

func (g *Creator) Create() error {
	return nil
}

func (g *Creator) setupPaths(args []string) error {
	return nil
}

func (g *Creator) setupConfig() error {
	return nil
}

func (g *Creator) runBeforeScripts() error {
	return nil
}

func (g *Creator) runAfterScripts() error {
	return nil
}

func (g *Creator) processFiles() error {
	return nil
}
