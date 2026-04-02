package generator

import (
	"errors"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"proj/internal/paths"
)

var (
	ErrNotInstalled    = errors.New("template not installed")
	ErrNotFromGit      = errors.New("template was not installed from git")
	ErrHasLocalChanges = errors.New("template has uncommitted local changes")
)

func init() {
	_ = ErrNotInstalled
}

type Uninstaller struct {
	cfg   *Config
	paths *paths.Paths
}

func NewUninstaller(cfg *Config) (*Uninstaller, error) {
	return &Uninstaller{
		cfg:   cfg,
		paths: cfg.Paths,
	}, nil
}

func (u *Uninstaller) Uninstall() error {
	targetPath := u.paths.TemplatePath

	if _, err := os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
		slog.Error("Template not found", slog.String("path", targetPath))
		return ErrNotInstalled
	} else if err != nil {
		slog.Error("Failed to check template path", slog.Any("error", err))
		return err
	}

	repo, err := git.PlainOpen(targetPath)
	if err != nil {
		if u.cfg.Force {
			slog.Warn("Template not a git repo, forcing removal", slog.String("path", targetPath))
		} else {
			slog.Error("Template is not a git repository", slog.String("path", targetPath))
			return ErrNotFromGit
		}
	} else {
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

		if isDirty {
			if u.cfg.Force {
				slog.Warn("Template has local changes, forcing removal", slog.String("path", targetPath))
			} else {
				slog.Error("Template has uncommitted local changes", slog.String("path", targetPath))
				return ErrHasLocalChanges
			}
		}
	}

	if u.cfg.NoWrite {
		slog.Info("Dry run - would remove template", slog.String("path", targetPath))
		return nil
	}

	slog.Debug("Removing template", slog.String("path", targetPath))

	if err := os.RemoveAll(targetPath); err != nil {
		slog.Error("Failed to remove template", slog.Any("error", err))
		return err
	}

	slog.Info("Template uninstalled successfully", slog.String("path", targetPath))
	return nil
}
