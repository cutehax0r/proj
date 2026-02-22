proj(7) -- Project template system overview
===========================================

## NAME

proj - Overview of the proj template system

## DESCRIPTION

**proj** is a project template system that allows you to define reusable
templates for programming projects. This manual provides an overview of the
template system, variable resolution, and script execution.

## TEMPLATE SYSTEM

Templates are directories containing files that can be processed using Go's
text/template package. Template files use the standard Go template syntax
with double braces:

    {{.VariableName}}

### Built-in Functions

proj supports all standard Go template functions, plus:

## SEE ALSO

proj(1), proj.yml(5), proj-scripts(7), proj-lua(7), proj-variables(7), proj-template(7)
