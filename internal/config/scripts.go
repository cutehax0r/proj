package config

import (
	"errors"
	"log/slog"
	"os"
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

func NewScriptSpec(paths *paths.Paths) (ScriptSpec, error) {
	var result ScriptSpec
	var ok bool

	// target level
	path := strings.Join([]string{"definitions", viper.GetString("definition-name"), "scripts"}, ".")
	scripts := viper.Get(path)
	if scripts != nil {
		scriptMap, ok := scripts.(map[string]any)
		if !ok {
			slog.Error("ScriptSpec parse failure", slog.Any("Scripts", scripts))
			return result, errors.New("ScriptSpec parse failure: entire scripts section is busted in definition")
		}
		if before, ok := scriptMap["definition-before"].(string); ok {
			fp, err := paths.ResolveTemplate(before)
			if err != nil {
				slog.Error("Bad path for definition-before", slog.String("definition-before", before), slog.Any("err", err))
				return result, err
			}
			result.DefinitionBefore = fp
		}
		if after, ok := scriptMap["definition-after"].(string); ok {
			fp, err := paths.ResolveTemplate(after)
			if err != nil {
				slog.Error("Bad path for definition-after", slog.String("definition-after", after), slog.Any("err", err))
				return result, err
			}
			result.DefinitionAfter = fp
		}
	}

	csm := viper.GetStringMapString("scripts")

	// template level
	if result.TemplateBefore, ok = csm["template-before"]; ok {
		fp, err := paths.ResolveTemplate(result.TemplateBefore)
		if err != nil {
			slog.Error("Bad path for template-before")
			return result, err
		}
		result.TemplateBefore = fp
	}
	if result.TemplateAfter, ok = csm["template-after"]; ok {
		fp, err := paths.ResolveTemplate(result.TemplateAfter)
		if err != nil {
			slog.Error("Bad path for template-After")
			return result, err
		}
		result.TemplateAfter = fp
	}

	// global level
	if result.GlobalBefore, ok = csm["new-before"]; ok {
		fp, err := paths.ResolveGlobal(result.GlobalBefore)
		if err != nil {
			slog.Error("Bad path for new-before")
			return result, err
		}
		result.GlobalBefore = fp
	}
	if result.GlobalAfter, ok = csm["new-after"]; ok {
		fp, err := paths.ResolveGlobal(result.GlobalAfter)
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
	result = appendIfExists(result, s.GlobalBefore)
	result = appendIfExists(result, s.TemplateBefore)
	result = appendIfExists(result, s.DefinitionBefore)
	return result
}

func (s *ScriptSpec) AfterScripts() []string {
	var result []string
	result = appendIfExists(result, s.DefinitionAfter)
	result = appendIfExists(result, s.TemplateAfter)
	result = appendIfExists(result, s.GlobalAfter)
	return result
}

func appendIfExists(scripts []string, script string) []string {
	if script == "" {
		return scripts
	}
	_, err := os.Stat(script)
	if err != nil {
		slog.Debug("Script doesn't exist", slog.Any("error", err), slog.String("script", script))
		return scripts
	}
	return append(scripts, script)
}
