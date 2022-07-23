---
related-docs: [milpa/command/environment, milpa/repo/index, milpa/repo/docs]
related-commands: ["itself create", "itself docs generate"]
index: Command specs
---
Command specs go along with your scripts and help inform `milpa` of what its input should look like. Based on it, `milpa` will produce help pages and autocompletions, and may validate the arguments to your command.

Specs must be written in YAML, have a `yaml` extension, and be named exactly like your command (minus the extension, if any). For example, given a command at `.milpa/commands/my-command.sh` the corresponding spec file would be `.milpa/commands/my-command.yaml`.

## Example

If we wanted to have a command written in bash at `milpa release`, we'll need to create file at `.milpa/commands/release.sh` and its corresponding spec at `.milpa/commands/release.yaml`. Let's take a look at what such a spec might look like:

```yaml
# a summary is required. It shows up during autocomplete and command listings
summary: Create a github release
# a description is required as well. Add a longer description of how your command does its magic here
description: |
  This command shows you a changelog and waits for approval before generating and pushing a new tag, creating a github release, and opening the browser at the new release.

  ## Schemes

  You may choose between [semver](https://semver.org) and [calver](https://calver.org). Their composition is as follows:

  - **semver**: `{major}.{minor}.{patch}`, i.e. `3.14.1592`
  - **calver**: `{date}.{minor}.{micro}`, where date is derived from the `prefix` option; for example `16.18.339`.

# arguments is an ordered list of to arguments expected by your command
arguments:
  # each argument gets a name, your script will find these in the environment
  # so in this example, increment becomes $MILPA_ARG_INCREMENT
  # but it's also available as the first argument (i.e. `$1` in bash)
  # `!$\/%^@#?:'"` are all disallowed characters in argument names.
  - name: increment
    # arguments require a description
    descrption: the increment to apply to the last git version
    # arguments can have a default
    default: patch
    # or be required, but not have a default and be required at the same time.
    required: true
    # arguments may be variadic, that is, be all the arguments passed by the user
    # starting at the position of this argument forwards, in this case, since
    # there's only one argument, it would mean all arguments after the command name
    # (not including options)
    variadic: true
    # argument values can come from many sources:
    # dirs, files, milpa, script or static can be used for any argument
    # if specified, arguments will be validated by default before running the command
    values:
      # autocompletes directory names only, if a prefix is passed
      # then it'll be used as a prefix to search from
      dirs: "prefix"
      # autocompletes files with the given extensions, if any
      files: [yaml, json, hcl]
      # milpa runs the subcommand and offers an option for every line of stdout
      # Options and arguments may be used within go templates. For example,
      # the following would execute `milpa itself increments --scheme semver`
      # {{ Current }} will insert the argument's current value durint autocompletion
      milpa: itself increments {{ Opt "scheme" }}
      # script runs the provided command with `bash -c "$script"` and offers
      # each line of stdout as an option during autocomplete
      # go templates can be used in `script` as well
      script: git tag -l {{ Opt "prefix" }}
      # arguments can be validated to be part of a static set of values
      static: [micro, patch, minor, major]
      # when using script or milpa values, wait at most this
      # amount of seconds before giving up during autocomplete or validation
      timeout: 10
      # only suggest these as autocompletions but don't validate them before running
      suggest-only: true
      # if enabled, will not add a space after a suggestion
      suggest-raw: false


# options, also known as flags, are specified as a map
options:
  # this sets the --scheme option
  # it will be available to your script as the $MILPA_OPT_SCHEME environment variable
  # your users will be able to use `--scheme "semver"` or `--scheme=semver` for example
  scheme:
    # options require a description
    description: Determines the format of the tags for this repo.
    # Sometimes, very commonly used flags might benefit from setting a short name
    # in this case, users would be able to use `-s calver`
    short-name: s
    # and they may have defaults
    default: semver
    # or be required, but not have a default and be required at the same time.
    required: true
    # as well as arguments, flags can be specified multiple times and be passed as a list to your script
    repeated: false
    # option values can come from/be validated against many sources as well!
    values:
      # autocompletes directory names only, if a prefix is passed
      # then it'll be used as a prefix to search from
      dirs: "prefix"
      # autocompletes files with the given extensions, if any
      files: [yaml, json, hcl]
      # milpa runs the subcommand and returns an option for every line of stdout
      milpa: "itself repo list"
      # script runs the provided command with `bash -c "$script"` and returns an
      # option for every line of stdout
      script: "git tag -l"
      # arguments can be validated to be part of a static set of values
      static: [micro, patch, minor, major]
      # when using script or milpa values, wait at most this
      # amount of seconds before giving up during autocomplete or validation
      timeout: 10
      # only suggest these as autocompletions but don't validate them before running
      suggest-only: true
      # if enabled, will not add a space after a suggestion
      suggest-raw: false

  prefix: # becomes MILPA_OPT_PREFIX
    description: |
      An optional prefix to prepend to release identifiers. If `calver` is chosen as `scheme`, you may specify a combination of `YY`, `YYYY`, `MM`, and `DD` to be replaced with the corresponding values of the local date. The default in that case is `YY`.
    # if desired, defaults can be set to empty strings to be explicit
    default: ""
  ask: # becomes MILPA_OPT_ASK
    description: Prompt for confirmation before creating a github release
    # flags can be boolean. Since environment variables can only be strings,
    # these false values (or defaults) will be set as an empty string, while
    # true values will be set as "true"
    type: bool
    default: false

```
