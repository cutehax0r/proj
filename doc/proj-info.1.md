proj-info(1) -- Inspect templates and definitions
===============================================

## SYNOPSIS

`proj info`

`proj info` [*options*]

`proj info` <template>

`proj info` <template> <definition>

## DESCRIPTION

Show information about templates and definitions.

With no arguments, list available templates.

With a template argument, show definitions grouped by what can be used for
project creation (`new`) and what can be used in existing projects (`add`).

With a template and definition argument, show details for that definition,
including target files and required variables.

When run inside a project directory (one containing `.proj/proj.yml`),
`proj info` shows the current template and local project definitions.

Use `--all` to force a global view and list all available templates,
even when inside a project directory.

## OPTIONS

* `-a, --all`:
  Force global view. List all available templates regardless of
  whether currently inside a project directory.

* `-s, --template-root <path>`:
  Path containing project templates.

## EXAMPLES

List available templates:

    proj info

Show definitions in a template:

    proj info static

Show details for one definition:

    proj info static new

Force global view when inside a project:

    proj info --all
    proj info -a

## SEE ALSO

proj(1), proj-new(1), proj-add(1), proj-template(7)
