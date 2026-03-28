package generator

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"proj/internal/paths"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

type Installer struct {
	cfg   *Config
	paths *paths.Paths
}

func NewInstaller(cfg *Config, templateRoot, templateGit string) (*Installer, error) {
	viper.Set("template-root", templateRoot)
	viper.Set("template-name", cfg.TargetName)

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return &Installer{
		cfg:   cfg,
		paths: cfg.Paths,
	}, nil
}

func (i *Installer) Install() error {
	templateGit := viper.GetString("template-git")
	cloneURL, err := i.buildCloneURL(i.cfg.TemplateName, templateGit)
	if err != nil {
		slog.Error("Failed to build clone URL", slog.Any("error", err))
		return err
	}

	i.cfg.DefinitionName = cloneURL

	targetPath := filepath.Join(i.paths.TemplateRoot, i.cfg.TargetName)

	if _, err := os.Stat(targetPath); err == nil {
		slog.Error("Template already exists", slog.String("target", i.cfg.TargetName), slog.String("path", targetPath))
		return errors.New("template already exists")
	}

	if i.cfg.NoWrite {
		slog.Info("Dry run - would clone template", slog.String("url", cloneURL), slog.String("to", targetPath))
		return nil
	}

	slog.Debug("Cloning template", slog.String("url", cloneURL), slog.String("to", targetPath))

	_, err = git.PlainClone(targetPath, false, &git.CloneOptions{
		URL:   cloneURL,
		Depth: 1,
	})
	if err != nil {
		slog.Error("Failed to clone", slog.Any("error", err))
		return err
	}

	projYmlPath := filepath.Join(targetPath, paths.TemplateConfigFile)
	if _, err := os.Stat(projYmlPath); err != nil {
		slog.Error("Template missing proj.yml", slog.String("path", projYmlPath))
		return err
	}

	slog.Debug("Template installed successfully", slog.String("path", targetPath))
	return nil
}

func (i *Installer) buildCloneURL(source, templateGit string) (string, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "git@") {
		return source, nil
	}

	parsed, err := url.Parse(templateGit)
	if err != nil {
		slog.Error("Invalid template-git URL", slog.Any("error", err))
		return "", err
	}

	owner := strings.Trim(parsed.Path, "/")

	return strings.TrimRight(parsed.Scheme+"://"+parsed.Host+"/"+filepath.Join(owner, source), "/"), nil
}
