
proj-update(1) -- Update installed templates
==============================================

## SYNOPSIS

`proj update` [target] [options]

## DESCRIPTION

Update a template currently installed.

Target is the name of a template installed in the template-root directory. If target is provided
then only that template will be updated.

If target is not specified all templates will be mapped over and updated if possible. Templates that
are not fetched from a git-source will be skipped.

## OPTIONS

* `-s, --template-root <path>`:
  Path containing project templates (default: ~/.local/share/proj)

* `-f, --force`:
  Delete and reinstall even if the repository has uncommitted changes

## EXAMPLES

Update all templates from GitHub (default):

    proj update

Update a specific template

    proj update html

## SEE ALSO

proj(1)
