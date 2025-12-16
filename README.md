# Proj
A tool for setting up new projects or adding files to existing projects

## Usage

Use the ruby on rails project template in `~/.local/proj/rails/new` to create a project called
`foo` in the current directory.

```sh
proj new rails foo
```

If you're in a rails project then this will add `app/models/foo.rb` and `spec/models/foo_spec.rb`
using the templates in `~/.local/proj/rails/files/model`

```sh
proj add model foo
```


## Development

General tips:

* Use `cobra-cli` to add commands
* Run with `go run main.go`
* Install modules with `go mod tidy`
* compile with `go build`
