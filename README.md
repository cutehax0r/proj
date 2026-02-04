# Proj
A tool for setting up new projects or adding files to existing projects

## Usage

Commands:
  * proj new <template> <name> [OPTIONS]

Running `new` implies the definition `new`

# Wishlist

* Break out variables from requirements into it's own thing
* Paths should gain 'resolve' methods for template, target, and global config paths
* Everything should switch from filepath.join to the new resolve methods
* Refactor `NewPaths()` to use a method rather than the if chain
* Config.go should have init methods for scripts, vars, etc. to simplify new

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
