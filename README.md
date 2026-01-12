# Proj
A tool for setting up new projects or adding files to existing projects

## Usage

Commands:
  * proj new <template> <name> [OPTIONS]
  * proj add <definition> <name> [OPTIONS]

Running `new` implies the definition `new`
Running `add` implies reading template from `.proj/proj.yml`

### Flags
* *dry-run*: don't actually change the filesystem, just log all the state (implies --log-level 2)
* *log-level*: amount of logging to display (from 0-3)
* *global-config*: specify a global config for use instead of the default
* *target-root*: location to create the project folder
* *target-path*: full path to where to create files (don't infer from target-path/name)
* *template-root*: path to where templates can be found
* *template-path*: specific tempate path to be used (don't infer from template-root/template)
* *set*: set a variable with key value pairs (--set variable=value). Can be passed multiple times

## Config
What is the minimal config that can work?

### Global
```yaml
---
# ~/.config/proj/proj.yml
template-root: ~/.local/share/proj
target-root: .
log-level: 0
scripts:
  new-before: scripts/before.lua
  new-after: scripts/after.lua
  add-before: scripts/before.lua
  add-after: scripts/after.lua
variables:
  global: #reserving space for variables.template && variables.template.definition
    aboolean: true
    anumber: 8080
    astring: hello world
    anarray: [1, 3, 3]
    amap:
      foo: 1
      bar: 2
      nested:
        yes: true
        no: false
        anotherarray:
          - 4
          - 5
    anil: ~
```

### Template
```yaml
---
# ~/.local/share/proj/html/proj.yml
scripts:
  template-before: scripts/before.lua
  template-after: scripts/after.lua
definitions:
  new:
    requirements:
      local: false
      variables:
        - name: "title"
        - name: "has_date"
        - name: "font_size"
          default: 14
        - name: "background"
          default: "#202020"
    scripts:
      definition-before: scripts/before.lua
      definition-after: scripts/after.lua
    files:
      - source: src/Dockerfile
        target: Dockerfile
        parse: true
      - source: src/Makefile
        target: Makefile
        parse: true
      - source: src/README.md
        target: README.md
        parse: true
      - source: src/index.html
        target: src/index.html
        parse: true
      - source: src/favicon.svg
        target: src/images/favicon.svg
        parse: false

  html:
    requirements:
      local: true
      variables:
        - name: "title"
    files:
      - source: src/index.html
        target: src/{{ .name }}.html
        parse: true

  css:
    requirements:
      local: true
    files:
      - source: src/index.css
        target: src/{{ .name }}.css
        parse: false

  js:
    requirements:
      local: true
    files:
      - source: src/index.js
        target: src/{{ .name }}.js
        parse: false
```

### Project
```yaml
---
# ~/src/my_homepage/.proj/proj.yml
kind: "html" # copied from when 'newed'
template_path: "~/.local/share/proj/html" # only if explicitly set when 'newed'
variables: # values that were explicitly set when `newd` otherwise we fall back to template/global
  global: # reserving space for variables.template and variables.definition
    title: foo
    fontsize: 12
    background: "#000000"
```

## Commands
All commands start by reading `~/.config/proj.yml` as `global_config`

### New
```sh
proj new html my_homepage --var "title=hello world"
```

* Do `viper.Set("command", "new")`
* Do `viper.Set("name", "my_homepage")`
* Do `viper.Set("definition", "new")`

* Ensure `.proj/proj.yml` does not exist: finding one is a fail.
  * Set `project_config_path` to `filepath.join(viper.Get("target_path"), ".proj", "proj.yml")` 
  * Set `project_config` to an empty `map[string]any`

* Do `viper.Set("kind", "html")` so `viper.Get("template_path")` is `~/.local/share/proj/html`

* Read `filepath.join(viper.Get("template_path"), "proj.yml")` as `template_config`
* Do `viper.Set("target_path, filepath.join(viper.Get("target_root"), name))`
* Ensure `viper.Get("target_path")` doesn't exist: finding one is a fail.

### Add

```sh
proj add html about_me --var "title=hello world"
```

* add won't take a target-root target-path or template root/template path. They are looked up from targetconfigfile
* add will have to set template name based on what it rad from the project config file

* Do `viper.Set("command", "add")`
* Do `viper.set("name", "about_me")`
* Do `viper.set("definition", "html")`

* Ensure `.proj/proj.yml` exists in directory hierarchy: Not finding one is a fail.
  * Set `viper.Set("project_config_path", project_config_path)`
* Read `viper.Get("project_config_path")` as `project_config` and
  * `viper.Set("kind", project_config['kind'])`
  * `viper.Set("template_root", project_config['template_path'])` if present
* Read `filepath.Join(viper.Get("template_path"), "proj.yml")` as `template_config`
* Do `viper.Set("target_path", filepath.Dir(filepath.Dir(Viper.Get("project_config_path"))))`

## Variable resolution

1. `resolved = map[string]any`
2. for each `variable` in `template_config['definitions.definition.requirements.variables']`
   `resolved[variable.name]` = `variable.default`
3. for each `variable` in `global_config['variables.global']`
   `resolved[variable.name]` = `variable.value`
4. for each `variable` in `project_config['variables.global'] -> name, value`
   `resolved[variable.name]` = `variable.value`
5. parse *argvars* for each `viper.GetStringSlice("set") -> (key, value)`
   `resolved[variables.key]` = `global_config['key']` (because viper will have set this form the arg)

* variables are set in config files. Nothing special happening here.
  * I'm choosing not to template eval the whole config file because I don't want programmign to
  happen both in lua and yaml.  If you need code stuff like setting a default for one thing based on
  another then do it in lua.
  * this means I am choosing to run go template resolution on each 'target' in the definition
  files declaration.
  * this means variable resolution must happen before any file copying can happen.

* variables are set/changed by lua scripts. I'll write a `proj.exec()` function that lets you run a shell
  command to do something like load a secret from keychain.

We're going to need to figure out where flags come from in order to handle the 'overriding' that has
to happen to make this all feel good.

The priority system is:

1. Args. Args always win.
2. Project config. If this exists then it carries an explicit change from the defaults or info we
   gathered at one time in the past and don't want to nag for in the future
3. global config. If set I'm assuming this was intentional to replace any template-level defaults
4. template config default

In order to make this work I'm going to have to interrogate the `cmd.Flags()` and keep a list of any
`--set` calls that were made

```go
cmd.Flags().StringArray("set", []string{}, "set variables with key=value pairs")
viper.BindPFlag("set", cmd.Flags().Lookup("set"))
```
```go
argvars := viper.GetStringSlice("set") // []string{"Buddy=dog", "Mittens=cat", ...}
for _, s := range runSet {
  key, value, found := strings.Cut(s, "=")
  if !found { continue }
  viper.Set("variables." + key, value)
}
```

### Nested Maps
This resolution happens only at the top level of the `resolved_variables` map.  if project config
has default `vars = { animals: { dog: buddy, sport: basketball} }` and global defaults has
`vars = { animals: { cat: mittens, sport: nil }}` then `animals` will be
`animals{cat: mittens, sport: nil}` because global vars win and that's the end of it.
No nested merging on for variables that are maps.

```go
// figure out where a flag was set
type FlagSource int

const (
  FlagSourceDefault FlagSource = iota
  FlagSourceConfig
  FlagSourceCLI
)

func flagSource(cfg *viper.Viper, cmd *cobra.Command, name string) FlagSource {
  flag := cmd.Flags().Lookup(name)
  switch {
    case flag != nil && flag.Changed:
      return FlagSourceCLI
    case cfg.InConfig(name):
      return FlagSourceConfig
    default:
      return FlagSourceDefault
  }
}
```

## Execution order

1. config files are read
  * global_config `~/.config/proj.yml`
  * project_config `./.proj/proj.yml`
  * template_config `~/.local/share/proj/kind/proj.yml`

2. resolve variables based on above process
3. Evaluate "before" scripts:
    1. global before `global_config.Get("scripts." + viper.Get("action") + "_before")` if present
    2. template global before `template_config.Get("scripts." + "template_before")` if present
    3. template definition before `template_config.Get(viper.get("definition") + ".scripts." + "definition_before")` if present

4. ensure resolved variables cover all requirements: for each `reqdvar` in
   `template_config[definitions.definition.requirements.variables']`
   if nil, prompt for value until not nil using options, type, to drive a UI.

5. for each file in `template_config['definitions.' + viper.Get("definition") '.files']` as file do
  * `source` = `filepath.Join(viper.Get("template_root"), file['source'])`
  * `target` = `GO_TEMPLATE_STRING(filepath.Join(viper.Get("target_path"), file['target'])), RESOLVED_VARIABLES)`
  * `content` = `GO_TEMPLATE_FILE(source, RESOLVED_VARIABLES)`
  * `file.Write(target, content)`

6. run after scripts
    7. template definition after `template_config.Get(viper.get("definition") + ".scripts." + "definition_after")` if present
    8. template global after `template_config.Get("scripts." + "template_after")` if present
    9. global after `global_config.Get("scripts." + viper.Get("action") + "_after")` if present

## Module structure
* main (entry point) - just runs commands
* cmd - all the commands exposed via cobra
* lua_runtime - code for setting up and enhancing the lua runtime
* config - handles loading vipers and preparing the variables for exection

# Wishlist

*GET IT WORKING BEFORE FUSSING WITH THIS STUFF*

* allow for `variables.template.definition` in global config so that (for example)
`variables.python.version = 3` or `variables.ruby.typesystem = "sorbet"` can be set
* allow variable declarations to include "persist: false" to avoid adding them to a .proj directory
* allow required variables to declare kind: (string, number, boolean) as type
* allow required variables to declare "options"
* write a `prompt` function that can be called from lua
* write a `shell exec` function exposed to lua that can run commands and recieve std err/in
* add "descriptions" to most things
* allow requirements to define mandatory software, checked with `which foo`
* add --action param on new that lets you have multiple 'new' types and you can choose between them
* add `setup` command that writes a default config to and an example template
* add `info` command that shows details about what's allowed/required for templates/definitions
* add `install` command that does a `git clone ...` to template_root
* add `remove` command that does an `rm -rf ` in template root
* make `new` list all available templates if run with no arguments
* make `add` list all available definitions if run with no arguments
* make `dry run` actually do nothing
* make a prettier logger
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
