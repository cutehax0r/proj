# Proj v0.0.14
A tool for setting up new projects or adding files to existing projects

## Usage

Commands:
  * `proj new <template> <name> [OPTIONS]`
  * `proj add <definition> <name> [OPTIONS]`
  * `proj info [template] [definition] [OPTIONS]`

Running `new` implies the definition `new`

### Command Line Arguments

  * `-l INT`, `--log-level INT`: How much to log. 0 = less, 3 = more

  * `-w BOOL`, `--no-write BOOL`: Wether to 'dry run' (true) or write changes to disk (false)

  * `-g STRING`, `--global-config-file STRING`: Specify an alternative config file

### Commands


#### New

Create a new project using a template

##### Usage

General format: `proj new <template> <name> [OPTIONS]`.

**Example:** Create a new project. This uses the `new` definition in the `html` template. The
variable `title` will be set to "My Home Page". The project will be created in the current working
directory at `./my_homepage`.

```sh
proj new html my_homepage -v "title=my home page"
```

##### Arguments

  * `-r STRING`, `--target-root STRING`: directory that will contain the new project folder

  * `-p STRING`, `--target-path STRING`: full path to the new project directory. Root to where files
    will be copied. Setting `--target-path` will cause `--target-root` to have no effect.

  * `-s STRING`, `--template-root STRING`: directory to search for project templates

  * `-t STRING`, `--template-path STRING`: full path for template that will be used to create new
    project. Setting `--template-path` will cause `--template-root` to have no effect.

  * `-v STRING`, `--set-varaible STRING`: define a variable `-v FOO=bar`.

  * `-d STRING`, `--definition-name STRING`: specific definition within the template to use

#### Add

Add a files to an existing project. May only be used with an existing project. You must either be
inside of a project directory or specify the `target-path`.

##### Usage

General format: `proj add <definition> <name> [OPTIONS]`

These examples assume an HTML template has been installed with reasonable definitions.

**Example:** Add an page to an existing project. Uses the `page` definition in the `html` template.
The varaible `title` will be set to `Contact Information` and `email` will be `root@example.com`.
The file would be written at `./contact.html`.

```sh
proj add page contact -v "email=root@example.com" -v "title=Contact Information"
```

**Example:** Add an a CSS file to an existing project. Uses the `stylesheet` definition in the
`html` template. The stylesheet would be at `./css/print.css`

```sh
proj add stylesheet print
```

##### Arguments

  * `-s STRING`, `--template-root STRING`: directory to search for project templates

  * `-t STRING`, `--template-path STRING`: full path for template that will be used to create new
    project. Setting `--template-path` will cause `--template-root` to have no effect.

  * `-v STRING`, `--set-varaible STRING`: define a variable `-v FOO=bar`.

#### Info

Inspect available templates and definitions.

##### Usage

General formats:

* `proj info`

* `proj info <template>`

* `proj info <template> <definition>`

##### Examples

**Example:** List available templates.

```sh
proj info
```

**Example:** Show definitions in the `static` template grouped by `new` and `add`.

```sh
proj info static
```

**Example:** Show details for the `new` definition in `static`.

```sh
proj info static new
```

##### Arguments

  * `-s STRING`, `--template-root STRING`: directory to search for project templates

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

### Variable Naming Conventions

All variable names are automatically normalized to lowercase. This means:

  * `UserName`, `userName`, and `username` all refer to the same variable

  * In templates and Lua scripts, use lowercase names: `{{.username}}` instead of `{{.UserName}}`

  * When defining variables in configuration files or via CLI, you can use any casing, but it will
  be converted to lowercase

**Note:** A warning will be logged if a variable name contains mixed case (e.g., `UserName`) to help
identify potential issues.

The system automatically provides the following lowercase variables in all definitions:

  * `targetname` - The name of the project/file being created

  * `templatename` - The name of the template being used

  * `definitionname` - The name of the definition being executed

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

  * `logInfo`. log "optional" information. Visible at log level 2 or 3.

  * `logWarn`. log "warnings" information, errors that can be recovered or that will not hault
  processing. Visible at log level 1-3.

  * `logError`. log "errors" that should halt further processing. visible at all log levels

# Development

Use `make` to execute commands. There isn't a convenient way to compile manuals from typst pages so
I'm going to use `go-md2man`.

# Wishlist

## Short term goals

  * can get 'make release' pushing:

      * cutehax0r/homebrew-tap

      * cutehaxor/rpm-repo

      * cutehaxor/deb-repo

      * cutehaxor/zst-repo

  * Docs pushing to github wiki on release

## Long term goals

  * **Fix design flaw:** File pulled from github aren't goign to have permissions copied. This is a
    problem for installing templates distributed as git repos. For example, If a template has a shell
    script it's execute bit won't be set. To fix this I think we need to modify the `files` declaration
    to allow you to declare target permissions and then not have them set on local files.

  * improve readme, focus on 'what it does' and how to install. Everything else to docs

  * Much improved installation instructions (binaries 'sudo xattr -d com.apple.quarantine ./path/to/app')
    or brew install --no-quarantine. linux repos need some kind of code signing too but it's gpg

  * github pages with a nice video intro, prettified documentation, etc. Get it syncing with releases.

  * allow for `variables.template.definition` in global config so that (for example)
  `variables.python.version = 3` or `variables.ruby.typesystem = "sorbet"` can be set

  * allow variable declarations to include "persist: false" to avoid adding them to a .proj directory

  * allow a 'not required' variable declaration format. for documentation of what's allowed

  * allow required variables to declare kind: (string, number, boolean) as type

  * allow required variables to declare "options" (must be 1,2, "alpha" or "bravo", true/false, etc)

  * write a `prompt` functions that can be called from lua

  * add "descriptions" to most things

  * allow requirements to define mandatory software, checked with `which foo`

  * add `setup` command that writes a default config to and an example template

  * add `create` that builds a new template in ~/.local/share/proj

  * add `install` command that does a `git clone ...` to template_root

  * add `remove` command that does an `rm -rf ` in template root

  * make `new` list all available templates if run with no arguments

  * make `add` list all available definitions if run with no arguments

  * add --version that dumps a version

  * shell completion should be better. use ValidArgsFunc() on cobra.Command to read template root. 

  * add support for --input or something that reads in data from a file (or maybe STDin and exposes it
    to lua-land and perhaps as a variable. That would allow you to automate filling content of certain
    files.

  * add string helper functions `camelcase`, `classify`, etc. like [active
  support](https://apidock.com/rails/ActiveSupport/Inflector/camelize)

  * a --force-set / -V(?) that forces variables to a value (so it gets applied after all the
  scripts run)

  * a --force-write that ignores the normal 'file already exists' checks

  * allow file definitions to specify destination permissions

  * All this fussing with variables.  It might make more sense to route everything to JSON and pass
  that to the toluavalue function.  Tag your structs with `json:"foo"` pass an instance of the struct
  to json.Marshal(somestruct) to get a map[string]any. then pass that over to .toluavalue Investigate
  later. For now it works and that's good enough.  Just write some more ToMap() functions

  * when copying files with parent directories, we should have a way to specify permissions. For now
  those will have to be done with shell scripts.

  * github windows builds working, add a winget repository
