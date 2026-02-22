proj.yml(5) -- Project configuration file format
=================================================

## NAME

proj.yml - Project configuration file for proj

## SYNOPSIS

The `proj.yml` file is the main configuration file for a proj project. It
defines project metadata, variables, scripts, and template settings.

## DESCRIPTION

The `proj.yml` file is a YAML configuration file that sits at the root of a
proj project. It contains:

* Project metadata (name, version, description)
* Variable definitions
* Script hooks for pre/post processing
* Template configuration

## FILE FORMAT

The file is formatted in YAML. A minimal example:

```yaml
variables:
  author: John Doe
  license: MIT
```

### Top-level Fields

* `name` (string, required):
  The name of the project.

* `variables` (map):
  Key-value pairs defining variables for template processing.

* `scripts` (map):
  Script hooks for execution at various stages.

* `templates` (map):
  Template-specific configuration.

## EXAMPLES

A complete example:

```yaml
name: my-service
```

## SEE ALSO

proj(1), proj(7), proj-scripts(7), proj-template(7), proj-variables(7)
