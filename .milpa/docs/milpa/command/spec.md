---
title: milpa command spec file
related-docs: [milpa/command/environment, milpa/repo/index, milpa/repo/docs]
related-commands: ["itself create", "itself docs generate"]
---
Command specs go along with your scripts and help inform `!milpa!` of what its input should look like. Based on it, `!milpa!` will produce help pages and autocompletions, and may validate the arguments to your command.

Specs must be written in YAML, have a `yaml` extension, and be named exactly like your command (minus the extension, if any). For example, given a command at `.milpa/commands/my-command.sh` the corresponding spec file would be `.milpa/commands/my-command.yaml`.

## Example

If we wanted to have a command written in bash at `!milpa! release`, we'll need to create file at `.milpa/commands/release.sh` and its corresponding spec at `.milpa/commands/release.yaml`. Let's take a look at what such a spec might look like:

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
    # arguments can be validated to be part of a static set of values
    values: [micro, patch, minor, major]
    # commands may request values for validation dynamically by reading the
    # lines printed by running another !milpa! subcommand.
    # arguments may have `values` or `values-subcommand`, but not both.
    values-subcommand: scm versions
    # arguments can have a default
    default: patch
    # or be required, but not have a default and be required at the same time.
    required: true
    # arguments may be variadic, that is, be all the arguments passed by the user
    # starting at the position of this argument forwards, in this case, since
    # there's only one argument, it would mean all arguments after the command name
    # (not including options)
    variadic: true

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
    # options can also be validated against a set of values
    values: [semver, calver]
    # and they may have defaults
    default: semver
    # or be required, but not have a default and be required at the same time.
    required: true
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

# you may specify related commands, help topics, and websites that are related to this command
# these will be suggested to the user on help screens
see-also:
  # the user will be suggested to look at `!milpa! help scm versions`
  - scm versions
  # the user will be suggested to look at `!milpa! help docs sdlc releasing`
  - help docs sdlc releasing
  - https://docs.github.com/en/rest/reference/repos#create-a-release
  - https://semver.org/
```
