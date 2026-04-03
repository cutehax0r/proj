package generator

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"proj/internal/paths"
)

type Updater struct {
	cfg   *Config
	paths *paths.Paths
}

func NewUpdater(cfg *Config) (*Updater, error) {
	return &Updater{
		cfg:   cfg,
		paths: cfg.Paths,
	}, nil
}

func (u *Updater) Update() error {
	templateRoot := u.paths.TemplateRoot
	templateName := u.cfg.TargetName

	targetPath := filepath.Join(templateRoot, templateName)

	if _, err := os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
		slog.Error("Template not found", slog.String("path", targetPath))
		return ErrNotInstalled
	} else if err != nil {
		slog.Error("Failed to check template path", slog.Any("error", err))
		return err
	}

	repo, err := git.PlainOpen(targetPath)
	if err != nil {
		slog.Info("Template not from git, skipping", slog.String("path", targetPath))
		return nil
	}

	remoteURL, err := u.getRemoteOrigin(repo)
	if err != nil {
		slog.Error("Failed to get remote origin", slog.Any("error", err))
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		slog.Error("Failed to get worktree", slog.Any("error", err))
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		slog.Error("Failed to get git status", slog.Any("error", err))
		return err
	}

	isDirty := !status.IsClean()

	if isDirty && !u.cfg.Force {
		slog.Info("Template has local changes, skipping", slog.String("path", targetPath))
		return nil
	}

	if isDirty && u.cfg.Force {
		slog.Warn("Template has local changes, force updating", slog.String("path", targetPath))
	}

	err = u.updateOrReinstall(repo, targetPath, remoteURL)
	if err != nil {
		return err
	}

	slog.Info("Template updated", slog.String("name", templateName), slog.String("path", targetPath))
	return nil
}

func (u *Updater) getRemoteOrigin(repo *git.Repository) (string, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return "", err
	}

	for _, r := range remotes {
		if r.Config().Name == "origin" {
			cfg := r.Config()
			if len(cfg.URLs) > 0 {
				return cfg.URLs[0], nil
			}
		}
	}

	return "", errors.New("no origin remote found")
}

func (u *Updater) updateOrReinstall(repo *git.Repository, targetPath, remoteURL string) error {
	forceReinstall := u.cfg.Force

	if !forceReinstall {
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		})
		if err != nil {
			if strings.Contains(err.Error(), "already up-to-date") || strings.Contains(err.Error(), "Already up to date") {
				return nil
			}
			slog.Warn("Pull failed, reinstalling", slog.String("path", targetPath), slog.Any("error", err))
			forceReinstall = true
		}
	}

	if forceReinstall {
		slog.Warn("Reinstalling template", slog.String("path", targetPath))
		return u.reinstall(targetPath, remoteURL)
	}

	return nil
}

func (u *Updater) reinstall(targetPath, remoteURL string) error {
	if u.cfg.NoWrite {
		slog.Info("Dry run - would reinstall template", slog.String("path", targetPath))
		return nil
	}

	if err := os.RemoveAll(targetPath); err != nil {
		slog.Error("Failed to remove template", slog.Any("error", err))
		return err
	}

	_, err := git.PlainClone(targetPath, false, &git.CloneOptions{
		URL:   remoteURL,
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

	return nil
}
