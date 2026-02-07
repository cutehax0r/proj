# Proj
A tool for setting up new projects or adding files to existing projects

## Usage

Commands:
  * proj new <template> <name> [OPTIONS]

Running `new` implies the definition `new`

## Configuration

### Global configuration

Global configuration options are set in the `XDG_CONFIG_DIR` or `~/.config/proj/proj.yml`.

  * Variables defined in the global configuration file will be available to all project templates or
    project-level custom definitions.

  * Scripts defined in the global configuration will before or after any template definition is run.
    You can define scripts that are run before 'add' or before 'new' commands.

  * Log level can be set from 0 to 3, with 0 being "log the least" (only errors) and 3 being "log
    the most" (debug-level).

  * Template root defaults to `XDG_DATA_DIR` or `~/.local/share/proj` and that is the location where
    template folders containing template defintions are expected. You can define a custom location
    which may be useful for development.

  * Target root sets the location in which new files are created. By default the current directory
    is used but you might want to set a `src` directory if you alawys want your projects to be
    started in a location (e.g., for a mono repo).
   
Here is an example configuration file:
```yaml
---
# ~/.config/proj/proj.yml
template-root: ~/src/github.com/cutehax0r/customtemplates
target-root: ~/src/github.com/cutehax0r/
log-level: 3
scripts:
  new-before: scripts/new-before.lua
  new-after: scripts/new-after.lua
  add-before: scripts/add-before.lua
  add-after: scripts/add-after.lua
variables:
  github:
    username: cutehax0r
    private: true
  sql:
    server: dev.myserver.local
    port: 3600
    user: root
    password: "password" # real value extract from password vault in new-before lua script
```

### Template Definition

Commands will read templates defined in `~/.local/share/proj`. The <template> argument in the
command selects a folder that must contain a `proj.yml`. An example template definition is:

#### Directory structucture
```sh
❯ ls -r
 ./
├── 󰣞 src/
│   ├── 󰕙 favicon.svg
│   ├──  index.html
├──  proj.yml
└── 󰂺 README.md
```

#### README.md
The readme should provide a description about the template. You may also which to include a
`license.md` if you're sharing the project template definition. These are currently unused but will
eventually be read by the `info` and `install` commands to provide additional details to the user.

#### Proj.yml example
This defines the 'new' method which creates a basic HTML project with index.html and favion.svg.
It also defines an 'add' method that can be used to add an HTML file to the project after it's been
created.
```yaml
---
# ~/.local/share/proj/html/proj.yml
description: "An minimal single page HTML website"
scripts: {}
  # template-before: scripts/html-template-before.lua
  # template-after: scripts/html-template-after.lua
definitions:
  # Create a new HTML project
  new:
    scripts: {}
      # definition-before: scripts/new-before.lua
      # definition-after: scripts/new-after.lua
    requirements:
      # only available if we're not already in a project
      # by default the `new` definition is run by `proj new` but you can have multiple new
      # definitions defined and call them with `proj new --definition foo`
      # by setting local false `proj new html foo` but inside an html project `proj add new` won't
      local: false # maybe this should have another name?
      # default is optional. Other variables may be defined in lua scripts but files will not be
      # written until all required variables are defined.
      variables:
        - name: "title"
          default: "Homepage"
    files:
      # source is relative to this file (proj.yml)
      # target is relative to the destination directory (./<name>)
      # target is a template string so you can use variables in ("foo/{{var}}.html")
      # if parse is false then template rendering is scoped and the file is copied as-is
      - source: "src/index.html"
        target: "index.html"
        parse: true
      - source: "src/favicon.svg"
        target: "favicon.svg"
        parse: false
  # Add an HTML file to an existing project
  html:
    scripts: {}
    requirements:
      # only available if we're are in a project. You can `proj add html foo`
      # but not `proj new html foo`
      local: false
      variables:
        # no default means this must be provided with -v title=foo
        - name: "title"
    files:
      # this uses the 'name' variable to specify the target file name
      - source: "src/index.html"
        target: "{{ name }}.html"
        parse: true
```

This structure is consistent.

```yaml
scripts: {}
variables: {}
definitions: {}
```

Within definitions you will have a series of

```yaml
defintion_name:
  scripts: {}
  requirements:
    local: false
    variables:
      - name: "port"
        default: 80
  files:
    - source:
      target:
      parse:
```

## Variable Resolution
By default variables set with command line arguments have the highest priority. `-v foo=bar` sets
the variable `foo` to the value `bar`.

If a variable is not set by the command line arguments then the template's definition is checked for
a default value. If no default value is set by the definition that is being run then the top level
of the template is checked for a variable declaration with a default value.

If no definition is set at the template level then the global configuration is checked for a
variable definition.

If any required variable is not set then processing will halt with an error message. Processing will
continue if a variable is undefined but not listed as required.

Once variables are declared they may be modified by any of the defined 'before scripts'. 

## Templates
Any file marked as `parse: true` will be ran through the standard go
[text/template](https://pkg.go.dev/text/template) processor.

`target` strings in file definitions within templates are also processed as templates and may use
variables.

## Scripts

Scripts are ran in order. Any declaration that happens in an 'earlier' script will be visible to any
script running afterwards. This means you can use the global `add-before` and `new-before` scripts
to declare modules, functions, or variables and all future scripts and templates will see the result
of that.

By default a `proj` module is defined. This module has the following attributes. These are intended
to be read only. While it is possible to modify them, the changes would only be visible to lua
scripts. Changes will not be propagated back to the template parser. Changing the value of `noWrite`
will not change the writing behavior

  * `noWrite`: Was the no-write flag passed at the command line. You should check this before
    running any code in script that might modify something permanently.

  * `paths`: A table containing the fully resolved paths for where to read files, where to write
    them, names of configuration files, etc.

  * `files`: the file sources/targets and parse states in the definition being run

  * `isLocal`: is the script running in the context of an `add` (false) or `new` (true) definition

  * `variables` a map containing the name and default value of each required variable.

A number of functions are also provided to help make scripts easier to write.

  * `logDebug`. Log "debug" information. Only visible at log-level 3.
  * `logInfo`. log "optional" information. Visible at log level 2 or 3.
  * `logWarn`. log "warnings" information, errors that can be recovered or that will not hault
  processing. Visible at log level 1-3.
  * `logError`. log "errors" that should halt further processing. visible at all log levels


# Wishlist

* allow for `variables.template.definition` in global config so that (for example)
`variables.python.version = 3` or `variables.ruby.typesystem = "sorbet"` can be set
* allow variable declarations to include "persist: false" to avoid adding them to a .proj directory
* allow required variables to declare kind: (string, number, boolean) as type
* allow required variables to declare "options"
* write a `prompt` function that can be called from lua
* write a `shell exec` function exposed to lua that can run commands and recieve std err/in
* add "descriptions" to most things
* allow requirements to define mandatory software, checked with `which foo`
* add `setup` command that writes a default config to and an example template
* add `create` that builds a new template in ~/.local/share/proj
* add `info` command that shows details about what's allowed/required for templates/definitions
* add `install` command that does a `git clone ...` to template_root
* add `remove` command that does an `rm -rf ` in template root
* make `new` list all available templates if run with no arguments
* make `add` list all available definitions if run with no arguments
* add --version that dumps a version
* makefile that does go build for arm/x86/apple.
* makefile that can build bundles (brew, arch, deb, rpm, nix)
* shell completion should be better. use ValidArgsFunc() on cobra.Command to read template root. 
* add support for --input or something that reads in data from a file (or maybe STDin and exposes it
  to lua-land and perhaps as a variable. That would allow you to automate filling content of certain
  files.

* document the 'overloading' that happens.
  * no-write implies log info or log debug
  * if you set template-path, then template root is 'ignored'
  * if you set target-path, then target root is ignored

* Bolt on context when returning errors  func foo() (int, error) { x, err := givesErr(); err != nil
  { return fmt.Errorf("addcontext %w", err) } }
* maybe --set-variables FOO=BAR should try to be case-insentive when matching variables?
* maybe we should have a --force-set that forces variables to a value (so it gets applied after all the
scripts run)

* All this fussing with variables.  It might make more sense to route everything to JSON and pass
that to the toluavalue function.  Tag your structs with `json:"foo"` pass an instance of the struct
to json.Marshal(somestruct) to get a map[string]any. then pass that over to .toluavalue Investigate
later. For now it works and that's good enough.  Just write some more ToMap() functions
