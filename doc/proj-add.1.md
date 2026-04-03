proj-add(1) -- Add files to a project template
==============================================

## SYNOPSIS

`proj add` <path> [options]

## DESCRIPTION

Add a file or directory to the current project template.

## OPTIONS

* `-p, --target-path <path>`:
  Directory containing .proj/proj.yml (defaults to current directory)

* `-s, --template-root <path>`:
  Path containing project templates (default: ~/.local/share/proj)

* `-t, --template-path <path>`:
  Path to read files from

* `-v, --set-variable <key=value>`:
  Set a variable using key=value

## EXAMPLES

Add a file to the current project:

    proj add src/myfile.txt

## SEE ALSO

proj(1)
