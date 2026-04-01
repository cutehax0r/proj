proj-uninstall(1) -- Uninstall a template
=========================================

## SYNOPSIS

`proj uninstall` <target> [options]

## DESCRIPTION

Remove a template that was installed from a git repository. This command
only works on templates that were installed via `proj install`, not manually
added templates.

The uninstaller will verify:
- The template exists in the template root
- The template was installed from git (not a manual addition)
- There are no uncommitted local changes

Use `--force` to bypass the git and local changes checks.

## OPTIONS

* `-s, --template-root <path>`:
  Path containing project templates (default: ~/.local/share/proj)

* `--force`:
  Force removal even if the template is not from git or has local changes.
  This will also override the `--no-write` flag if both are specified.

* `-w, --no-write`:
  Dry run mode. Show what would be removed without actually deleting files.

## EXAMPLES

Uninstall a template:

    proj uninstall my-template

Dry run to see what would be removed:

    proj uninstall my-template --no-write

Force removal without checking git status:

    proj uninstall my-template --force

## SEE ALSO

proj(1), proj-install(1)
