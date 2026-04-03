proj-new(1) -- Create a new project from a template
===================================================

## SYNOPSIS

`proj new` <template> [name] [options]

## DESCRIPTION

Create a new project from the specified template.

## OPTIONS

* `-d, --definition-name <name>`:
  Definition in template to use (default: "new")

* `-v, --set-variable <key=value>`:
  Set a variable using key=value

* `-p, --target-path <path>`:
  Path to write files at (defaults to pwd/name)

* `-r, --target-root <path>`:
  Path to create the project in (default: ".")

* `-s, --template-root <path>`:
  Path containing project templates (default: ~/.local/share/proj)

* `-t, --template-path <path>`:
  Path to read files from

## EXAMPLES

Create a new project from a template:

    proj new static my-website

## SEE ALSO

proj(1)
