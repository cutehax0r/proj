package generator

import "proj/internal/paths"

type Installer struct {
	cfg *Config
	paths *paths.Paths
}

func NewInstaller(cfg *Config) (*Installer, error) {
	return &Installer{
		cfg: cfg,
		paths: cfg.Paths,
	}, nil
}

func (i *Installer) Install() error {
	// arg has a 'name'
	// arg has a 'target'
	// config has a 'repo source'
	// set target to last part of 'name' if target not set
	// build url to thing to install. -- source+name || name if a url
	// ensure template-root/target doesn't exist - if not then error
	// git clone to template root
	// ensure template-root/name has proj.yml -- if not then error
	return nil
}
