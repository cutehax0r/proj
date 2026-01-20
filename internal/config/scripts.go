package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/paths"
	"strings"

	"github.com/spf13/viper"
)

type ScriptSpec struct {
	GlobalBefore     string
	GlobalAfter      string
	TemplateBefore   string
	TemplateAfter    string
	DefinitionBefore string
	DefinitionAfter  string
}

func ParseScriptSpecs(paths *paths.Paths) (ScriptSpec, error) {
	var result ScriptSpec

	// target level
	path := strings.Join([]string{"definitions", viper.GetString("definition-name"), "scripts"}, ".")
	scripts := viper.Get(path)
	scriptMap, ok := scripts.(map[string]any)
	if !ok {
		slog.Error("ScriptSpec parse failure", slog.Any("Scripts", scripts))
		return result, errors.New("ScriptSpec parse failure: entire scripts section is busted in definition")
	}
	if before, ok := scriptMap["definition-before"].(string); ok {
		fp, err := filepath.Abs(filepath.Join(paths.TemplatePath, before))
		if err != nil {
			slog.Error("Bad path for definition-before", slog.String("definition-before", before), slog.Any("err", err))
			return result, err
		}
		result.DefinitionBefore = fp
	}
	if after, ok := scriptMap["definition-after"].(string); ok {
		fp, err := filepath.Abs(filepath.Join(paths.TemplatePath, after))
		if err != nil {
			slog.Error("Bad path for definition-after", slog.String("definition-after", after), slog.Any("err", err))
			return result, err
		}
		result.DefinitionAfter = fp
	}

	csm := viper.GetStringMapString("scripts")

	// template level
	if result.TemplateBefore, ok = csm["template-before"]; ok {
		fp, err := filepath.Abs(filepath.Join(paths.TemplatePath, result.TemplateBefore))
		if err != nil {
			slog.Error("Bad path for template-before")
			return result, err
		}
		result.TemplateBefore = fp
	}
	if result.TemplateAfter, ok = csm["template-after"]; ok {
		fp, err := filepath.Abs(filepath.Join(paths.TemplatePath, result.TemplateAfter))
		if err != nil {
			slog.Error("Bad path for template-After")
			return result, err
		}
		result.TemplateAfter = fp
	}

	// global level
	if result.GlobalBefore, ok = csm["new-before"]; ok {
		fp, err := filepath.Abs(filepath.Join(paths.GlobalConfigRoot, result.GlobalBefore))
		if err != nil {
			slog.Error("Bad path for new-before")
			return result, err
		}
		result.GlobalBefore = fp
	}
	if result.GlobalAfter, ok = csm["new-after"]; ok {
		fp, err := filepath.Abs(filepath.Join(paths.GlobalConfigRoot, result.GlobalAfter))
		if err != nil {
			slog.Error("Bad path for new-After")
			return result, err
		}
		result.GlobalAfter = fp
	}
	return result, nil
}

func (s *ScriptSpec) BeforeScripts() []string {
	var result []string
	if scriptExists(s.GlobalBefore) {
		result = append(result, s.GlobalBefore)
	}
	if scriptExists(s.TemplateBefore) {
		result = append(result, s.TemplateBefore)
	}
	if scriptExists(s.DefinitionBefore) {
		result = append(result, s.DefinitionBefore)
	}
	return result
}

func (s *ScriptSpec) AfterScripts() []string {
	var result []string
	if scriptExists(s.DefinitionAfter) {
		result = append(result, s.DefinitionAfter)
	}
	if scriptExists(s.TemplateAfter) {
		result = append(result, s.TemplateAfter)
	}
	if scriptExists(s.GlobalAfter) {
		result = append(result, s.GlobalAfter)
	}
	return result
}

func scriptExists(s string) bool {
	if s == "" {
		return false
	}
	_, err := os.Stat(s)
	if err != nil {
		slog.Debug("Checking script exists failed", slog.Any("error", err), "script", s)
		return false
	}
	return true
}
