proj(1) -- Project template manager
===================================

## SYNOPSIS

`proj` <command> [options]

`proj` [ `-h` | `--help` ]

## DESCRIPTION

**proj** is a project template manager that allows you to define templates for
programming projects. A template describes a directory structure and the files
contained within.

## COMMANDS

* `new` <template> [name]
  Create a new project from a template.

* `add` <path>
  Add a file or directory to the current project template.

* `info` [template] [definition]
  Inspect templates and definitions.

* `install` <source> [target]
  Install a template from a git repository.

* `uninstall` <target>
  Uninstall a template.

* `update` [target]
  Update installed templates.

* `completion` [command]
  Generate autocompletion scripts for shell

* `help` [command]
  Show help for a command.

## OPTIONS

* `-h`, `--help`:
  Show help message and exit.

* `-w, --no-write`:
  Print the plan and don't actually write any files.

* `-l, --log-level <level>`:
  How much to log [0-3], bigger = more (default: 0)

* `-g, --global-config-file <path>`:
  Use specific global configuration file.

* `--template-git <url>`:
  Default git source for templates (default: https://github.com/)

## FILES

* `proj.yml`
  Project configuration file. See `proj.yml(5)` for details.

* `.proj/`
  Directory containing project-specific configuration and templates.

## EXAMPLES

Create a new project from a template:

    proj new static my-website

Install a template from GitHub:

    proj install cutehax0r/my-template

Update installed templates:

    proj update

## SEE ALSO

proj-new(1), proj-add(1), proj-install(1), proj-uninstall(1), proj-update(1), proj-info(1), proj-help(1), proj.yml(5), proj(7)
