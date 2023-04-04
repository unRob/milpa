---
related-docs: [milpa/command/environment, milpa/repo/index, milpa/repo/docs]
related-commands: ["itself create", "itself docs generate"]
index: Command specs
---
Command specs go along with your scripts and help inform `milpa` of what its input should look like. Based on it, `milpa` will produce help pages and word completions, and may validate the arguments to your command.

Specs must be written in YAML, have a `yaml` extension, and be named exactly like your command (minus the extension, if any). For example, given a command at `.milpa/commands/my-command.sh` the corresponding spec file would be `.milpa/commands/my-command.yaml`.

## Example

If we wanted to have a command written in bash at `milpa release`, we'll need to create file at `.milpa/commands/release.sh` and its corresponding spec at `.milpa/commands/release.yaml`. Let's take a look at what such a spec might look like:

```yaml
summary: Create a github release
description: |
  This command shows you a changelog and waits for approval before generating and pushing a new tag, creating a github release, and opening the browser at the new release.

  ## Schemes

  You may choose between [semver](https://semver.org) and [calver](https://calver.org). Their composition is as follows:

  - **semver**: `{major}.{minor}.{patch}`, i.e. `3.14.1592`
  - **calver**: `{date}.{minor}.{micro}`, where date is derived from the `prefix` option; for example `16.18.339`.
arguments:
  - name: increment
    descrption: the increment to apply to the last git version
    default: patch
    required: true
    values:
      static: [micro, patch, minor, major]
options:
  scheme:
    description: Determines the format of the tags for this repo.
    short-name: s
    default: semver
    values:
      static: [semver, calver]
  prefix:
    description: |
      An optional prefix to prepend to release identifiers. If `calver` is chosen as `scheme`, you may specify a combination of `YY`, `YYYY`, `MM`, and `DD` to be replaced with the corresponding values of the local date. The default in that case is `YY`.
    default: ""
  ask:
    description: Prompt for confirmation before creating a github release
    type: bool
    default: false
```

## The Basics

```yaml
# a summary is required. It shows up during autocomplete and command listings
summary: Create a github release
# a description is required as well, and describes how your command does its magic
# it may be formatted with markdown
description: |
  This command shows you a changelog and waits for approval before generating and pushing a new tag, creating a github release, and opening the browser at the new release.

  ## Schemes

  You may choose between [semver](https://semver.org) and [calver](https://calver.org). Their composition is as follows:

  - **semver**: `{major}.{minor}.{patch}`, i.e. `3.14.1592`
  - **calver**: `{date}.{minor}.{micro}`, where date is derived from the `prefix` option; for example `16.18.339`.
# a list of arguments expected by your command, if any
# as well as a map of options
# see below for more details on options and arguments
arguments: []
options: {}
```

## Arguments

The `arguments` list describes the positional arguments that may be passed to a command. Arguments require a `name` and a `description`. The `name` of the argument will become available to commands through the environment variable named `MILPA_ARG_$NAME` where `$NAME` means the uppercased value of the `name` property. For example, an argument with `name: increment` will be available as your command's environment variable `MILPA_ARG_INCREMENT`. Arguments are passed as positional arguments to your command, i.e. `$1`, `$2`, and so on. Dashes (`-`) will be turned into underscores `_` (`name: my-argument` turns into `MILPA_ARG_MY_ARGUMENT`).

> ⚠️ Argument names need to be valid as environment variable names, and characters `!$\/%^@#?:'"` are not allowed.

```yaml
# arguments is an ordered list of to arguments expected by your command
arguments:
  # available in to the command as the MILPA_ARG_INCREMENT environment variable
  - name: increment
    # this description shows up during auto-completion and in the command's help page
    descrption: the increment to apply to the last git version
    # a default may be specified, it'll be passed to your command if none is provided
    default: patch
    # if marked as required, the command won't run unless this argument is provided
    # An error will result if the argument is both required and has a default set
    required: true
    # arguments may be variadic, that is, all remaining arguments starting at this position
    # in this case, since there's only one argument, it would mean all arguments after
    # the command name (not including options)
    # should an argument be variadic, it's `default:` then must also be a list!
    variadic: true
    # the `values` property specifies how to provide completions and perform validation on
    # the values provided at the command line
    values: {}
```

## Options

The `options` map describes the named options that may be passed to a command. Options require a `name` and a `description`. The `name` of the option will become available to commands through the environment variable named `MILPA_OPT_$NAME` where `$NAME` means the uppercased value of key for an option. For example, an option at the key `scheme` will be available as your command's environment variable `MILPA_OPT_SCHEME`.

> ⚠️ Options are not available as positional arguments to your command. The same character restrictions as arguments apply to options.

```yaml
# options, also known as flags, are specified as a map
options:
  # this creates the --scheme option
  # it will be available to your script as the $MILPA_OPT_SCHEME environment variable
  # and may be specified on the command line as either `--scheme "semver"` or `--scheme=semver`.
  scheme:
    # options require a description, this will show during completions
    # and on the command's help page
    description: Determines the format of the tags for this repo.
    # Sometimes, very commonly used flags might benefit from setting a short name
    # in this case, users would be able to use `-s calver`
    short-name: s
    # a default value may be passed to the command if none is provided by the user
    default: semver
    # the values provided at the command line
    # flags can be boolean. Since environment variables can only be strings,
    # false values (the default) will be passed as an empty string "", while
    # true values will be passed as the string "true"
    type: bool # or `string`
    # the `values` property specifies how to provide completions and perform validation on
    values: {}
```
---

## Value completion and validation

A `values` property may be specified for both arguments and options; `milpa` will provide completion and validation from the following sources:

- directories (with a matching prefix),
- files (with given extensions),
- stdout of any `milpa` sub-commands,
- stdout of any given bash script, and
- a static list of pre-defined values.

```yaml
    # if specified for any `argument` or `option`
    #
    # `dirs`, `files`, `milpa`, `script`, or `static`
    values:
      # autocompletes directory names only. If a prefix is set, it'll be used as a
      # prefix to filter matches. In this example, it would offer directories with a name
      # that begins with `config` as completion values.
      # Validation of provided values is never performed
      dirs: "config"
      # autocompletes files with the given extensions, if any
      # Validation of provided values is never performed
      files: [yaml, json, hcl]
      # `milpa` runs the named milpa sub-command, offering a completion for every line of stdout
      # Options and arguments to sub-commands can be provided as go templates. For example,
      # the following would run `milpa itself increments --scheme semver`
      # {{ Current }} will insert the argument's current value durint completion
      milpa: itself increments {{ Opt "scheme" }}
      # script runs the provided command with `bash -c "$script"` and offers
      # each line of stdout as an option during autocomplete
      # go templates can be used in `script` as well
      script: git tag -l {{ Opt "prefix" }}
      # values can be validated to be part of a given set of strings
      static: [micro, patch, minor, major]
      # when using `script` or `milpa` values, wait at most this amount of seconds
      # before erroring out during autocomplete or validation
      timeout: 10
      # only suggest values as completions but don't validate them before running
      # has no effect for `dirs` or `files` as these are always suggestions and never validated
      suggest-only: true
      # if enabled, will not add a space after suggestions during autocomplete
      suggest-raw: false
```

### Value completion script interpolation

[go-template](https://pkg.go.dev/text/template#hdr-Actions) tags may be used within `milpa` and `script` value completions to interpolate already supplied values. The following tags are available:

- `{{ Arg "name" }}`: the value (or default) for the named argument. A map of argument names to their string values is also available at `{{ index Args "name" }}`.
- `{{ Opt "name" }}`: the value (or default) for the named option. A map of argument names to their string values is also available at `{{ index Opts "name" }}`.
- `{{ Current }}`: the value currently being auto-completed. On a command line like `milpa song play joão⇥`, `{{ Current }}` would return `joão`.

Additionally, the following Go functions are available:

- `{{ contains "team" "ea" }}`: [`strings.Contains`](https://pkg.go.dev/strings#Contains)
- `{{ hasSuffix "content" "ent" }}`: [`strings.HasSuffix`](https://pkg.go.dev/strings#HasSuffix)
- `{{ hasPrefix "content" "con" }}`: [`strings.HasPrefix`](https://pkg.go.dev/strings#HasPrefix)
- `{{ replace "file.yaml" ".yaml" ".sh" }}`: [`strings.ReplaceAll`](https://pkg.go.dev/strings#ReplaceAll)
- `{{ toUpper "shout" }}`: [`strings.ToUpper`](https://pkg.go.dev/strings#ToUpper)
- `{{ toLower "whisper" }}`: [`strings.ToLower`](https://pkg.go.dev/strings#ToLower)
- `{{ trim " padded " }}`: [`strings.Trim`](https://pkg.go.dev/strings#Trim)
- `{{ trimSuffix "content" "con" }}`: [`strings.TrimSuffix`](https://pkg.go.dev/strings#TrimSuffix)
- `{{ trimPrefix "content" "ent" }}`: [`strings.TrimPrefix`](https://pkg.go.dev/strings#TrimPrefix)

---

## _Group Command_ Spec

A _group command_, that is, a command whose only job is to serve as a grouping mechanism for other commands, can also have a limited spec. These are useful to provide clearer documentation on what that set of commands are to be used for, as well as setting group-level **options** (see above).

Let's say we have a few commands under the group `milpa test`: `milpa test unit` and `milpa test integration`. A `summary`, `description` and `options` can be defined for this group by creating a spec at `.milpa/commands/test/_test.yaml`, that is, the name of the group (in this case `test`) prefixed by a single underscore (`_`) with.For example:

```yaml
# .milpa/commands/test/_test.yaml
summary: Commands for running tests
description: |
  Holds commands to run unit tests and integration tests.
# any options defined in a group command spec
# will be available to any immediate children of `test`
options:
  coverage:
    type: bool
    description: generate coverage reports
```

> ⚠️ `arguments` are not allowed as part of a _group command_ spec.
