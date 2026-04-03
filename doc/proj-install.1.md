proj-install(1) -- Install a template from a git repository
============================================================

## SYNOPSIS

`proj install` <source> [target] [options]

## DESCRIPTION

Install a template to the local template-root from a git repository.

The source can be:

* A short name (e.g., `user/repo`) which will be prefixed with the default git source (https://github.com/)
* A full URL (e.g., `https://gitlab.com/user/repo`) which will be used as-is

If target is not specified, it defaults to the last component of the source.

## OPTIONS

* `-s, --template-root <path>`:
  Path containing project templates (default: ~/.local/share/proj)

* `-g, --template-git <url>`:
  Default git source for short template names (default: https://github.com/)

## EXAMPLES

Install from GitHub (default):

    proj install cutehax0r/my-template

Install with custom name:

    proj install cutehax0r/my-template my-cool-template

Install from GitLab:

    proj install cutehax0r/my-template --template-git https://gitlab.com/

Install from specific URL:

    proj install https://gitlab.com/me/my-template my-template

## SEE ALSO

proj(1)
