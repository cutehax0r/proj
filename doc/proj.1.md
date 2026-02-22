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

* `new` <template> [path]
  Create a new project from a template.

* `add` <path>
  Add a file or directory to the current project template.

* `completion` [command]
  Generate autocompletion scripts for shell

* `help` [command]
  Show help for a command.

## OPTIONS

* `-h`, `--help`:
  Show help message and exit.

## FILES

* `proj.yml`
  Project configuration file. See `proj.yml(5)` for details.

* `.proj/`
  Directory containing project-specific configuration and templates.

## EXAMPLES

Create a new project from a template:

    proj new my-template my-project

## SEE ALSO

proj-new(1), proj-add(1), proj-help(1), proj.yml(5), proj(7)
