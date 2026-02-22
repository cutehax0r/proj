proj-variables(7) -- Variable definition and resolution
=======================================================

## NAME

proj-variables - Variable definition and resolution order

## DESCRIPTION

Variables are key-value pairs used in templates. They can be defined in
multiple locations and are resolved in a specific order.

## RESOLUTION ORDER

Variables are resolved in the following order (later sources override earlier):

## DEFINITION

Variables are defined in `proj.yml` files:

```yaml
variables:
  name: my-project
  version: "1.0.0"
  debug: true
  ports:
    - 8080
    - 8081
```

## TYPES

Variables support the following types:

* **string** - Text values
* **int** - Integer numbers
* **float** - Floating-point numbers
* **bool** - Boolean (true/false)
* **table** - Arrays and maps

## LIMITATIONS

* Variable names are always lowercased internally
* Variable names must be valid Lua identifiers (alphanumeric + underscore)
* Nested tables are supported but limited to 3 levels deep

## SEE ALSO

proj(1), proj.yml(5), proj-template(7), proj-scripts(7)
