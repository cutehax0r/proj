proj-scripts(7) -- Script execution order
=========================================

## NAME

proj-scripts - Script execution order and hooks

## DESCRIPTION

Scripts are executable files that run at specific points during project
creation. They allow dynamic content generation, validation, and setup.

## EXECUTION ORDER

1. global before
1. template before
1. definition before
1. definition after
1. template after
1. global after

## SEE ALSO

proj(1), proj.yml(5), proj-variables(7), proj-template(7), proj-lua(7)
