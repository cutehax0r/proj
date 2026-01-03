# Proj
A tool for setting up new projects or adding files to existing projects
Maybe rename this to "jorp" - it's less likely to confligt with other people's stuff

## Usage

Use the ruby on rails project template in `~/.local/share/proj/rails/new` to create a project called
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

## General planning

* read config from `~/.config/proj`, then from `$CWD/.proj`.
* have templates stored in `~/.local/share/proj/templates`
* `templates` have a structure: `~/.local.share/proj/templates/html`
    * `README.md` a readme file
    * `config.yml` variables etc. that come from config
    * `before.lua(?)` a script that gets run before anything happens
    * `after.lua(?)` a script that runs after everything happens
    * `files/` the directory that will be 'cloned' to create the new project
    * `add/` templates for adding individual files
* `add/` holds files that look just like `templates` but don't have the `add` folder
* `add/remove.lua` can be run to remove a file

when you run `proj new html foobar`

  0. check that `./foobar` doesn't exist: crash if it does
  1. config loads from `~/.config/proj/config.yml`
  2. config loads from `~/.local/share/proj/templates/html/config.yml`
     verify sane config - resolve missing items || crash
  3. run `~/.local/share/proj/templates/html/before.lua`
     crash if return is non-zero
  4. copy `~/.local/share/proj/templates/html/files` to `./foobar` (copy = process template)
  5. run `~/.local/share/proj/templates/html/after.lua`
     crash if return is non-zero
  6. make `~/foobar/.proj/`

when you run `proj add index`

  0. check that `./proj` exists. traverse up the directory stack looking for it. crash if missing
  1. config loads from `~/.config/proj/config.yml`
  2. config loads from `./proj/config.yml`
  3. config loads from `~/.local/share/proj/templates/html/config.yml` (pull html from proj config)
  4. config loads from `~/.local/share/proj/templates/html/add/index/config.yml` (html from .proj)
     verify sane config - resolve missing items || crash
  5. run `~/.local/share/proj/templates/html/add/index/before.lua`
     crash if return is non-zero
  6. need some kind of 'pre compute destination file names before copying and crash if conflict
  7. copy `~/.local/share/proj/templates/html/add/index/files` to `./foobar` (copy = process template)
  8. run `~/.local/share/proj/templates/html/add/index/after.lua`
     crash if return is non-zero

config order of override: `~/.config`, `~/.local/share/proj/templates/foo/config`, `./proj/config`

contents of config
  * variables declarations:
    * can be set as string, int bool, array, etc.
    * can come from ENV or from process
  * variable 'requirements': optional, mandatory
  * pre-run/post-run scripts
  * file mapping:
    * files/index.html to .proj/../src/index.html
    * files/index_test.py to .proj/../test/index_test.html

We need some kind of 'setup' function you can run: 
* creates a default config in ~/.config/proj
* makes sure ~/.local/share/proj exists

some kind of an "info" service would be nice. Display the details about an installed project:
requirements, the readme, what variables are used (and current values), etc.

