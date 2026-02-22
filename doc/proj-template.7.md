proj-template(7) -- Template syntax and functions
=================================================

## NAME

proj-template - Template syntax and available functions

## DESCRIPTION

Templates use Go's `text/template` package for variable substitution and
logic. All standard Go template syntax is supported.

## BASIC SYNTAX

```
{{.VariableName}}
{{.Nested.Value}}
```

## CONDITIONALS

```
{{if .Debug}}
  debug mode enabled
{{else}}
  debug mode disabled
{{end}}
```

## LOOPS

```
{{range .Ports}}
  - {{.}}
{{end}}
```

## STANDARD FUNCTIONS

All Go `text/template` built-in functions are available:

* `and`, `or`, `not` - Logical operators
* `eq`, `ne`, `lt`, `le`, `gt`, `ge` - Comparison
* `index` - Access array/map by index
* `len` - Length of string/array
* `printf` - Formatted output

## ADDITIONAL FUNCTIONS

proj provides additional formatting functions:

this stuff doesn't exist yet but will be similar to active-support's helpers in rails

## EXAMPLES

```
package {{snake .Name}}

const Version = "{{.Version}}"

{{if .Debug}}
var Debug = true
{{end}}
```

## SEE ALSO

proj(1), proj-variables(7), proj-scripts(7)
