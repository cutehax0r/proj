package generator

import (
	"fmt"
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
		return err
	}

	i.cfg.DefinitionName = cloneURL

	targetPath := i.paths.TemplateRoot

	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("template %q already exists at %s", i.cfg.TargetName, targetPath)
	}

	slog.Debug("Cloning template", slog.String("url", cloneURL), slog.String("to", targetPath))

	_, err = git.PlainClone(targetPath, false, &git.CloneOptions{
		URL:      cloneURL,
		Depth:    1,
		Progress: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to clone: %w", err)
	}

	projYmlPath := filepath.Join(targetPath, paths.TemplateConfigFile)
	if _, err := os.Stat(projYmlPath); err != nil {
		return fmt.Errorf("template missing %s at %s", paths.TemplateConfigFile, projYmlPath)
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
		return "", fmt.Errorf("invalid template-git URL: %w", err)
	}

	owner := strings.Trim(parsed.Path, "/")

	return fmt.Sprintf("%s://%s/%s", parsed.Scheme, parsed.Host, filepath.Join(owner, source)), nil
}
