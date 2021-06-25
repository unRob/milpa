# milpa

_milpa_ is an agricultural method that combines multiple crops in close proximity. `milpa` is a Bash script and tool to care for one's own garden of scripts.

You and your team write scripts and a little spec for each command. Use bash, or any other command, and `milpa` provides autocompletions, sub-commands and argument parsing+validation for you to focus on your scripts. For those in the know, it makes following these [Command Line Interface Guidelines](https://clig.dev/) easier.

> This repo has just been planted, and everything here is extremely experimental.

## Usage

`milpa [sub-command...] [options] [arguments]`

### Layout

Let's say you have this in your repo:

```yaml
.milpa/
  commands/
    cloud-provider/
      login
      login.yaml
    db/
      connect.sh
      connect.yaml
      list.sh
      list.yaml
    vpn/
      connect
      connect.yaml
    onboard.sh
    onboard.yaml
```

Then, `milpa` would allow you to run `milpa cloud-provider login` and `milpa onboard`, as well as `milpa db connect api --environment production --verbose`, or even `milpa help db list` and so on and so forth. You choose how to organize your **Scripts** under `.milpa/commands`, and `milpa` figures out the rest.

### Where `milpa` looks for scripts

By default, `milpa` will look at the current git repo's root, or if git is not available, in your current working directory. Additional repositories can be added as colon (`:`) delimited paths to the folders containing a `/.milpa` folder within. For example, a `MILPA_PATH=$HOME/code/my-repo:$HOME/.local/milpa` would add the `$HOME/code/my-repo` and `$HOME/.local/milpa` folders to the command search path. Commands with the same name will override any commands previously found in the `MILPA_PATH`.

### Scripts

As it comes for scripts, you can either:

- write a bash script, with an `.sh` extension, or
- write whatever language you want, as long as you can name your script without an extension and make it executable.

In the example [layout above](#layout), `.milpa/commands/cloud-provider/login` could be a binary, and `.milpa/commands/vpn/connect` an applescript. `milpa` doesn't care as long as they have the executable bit set (i.e. don't forget to `chmod +x .milpa/commands/../your-script`!).

In order for `milpa` to recognize your commands, you'll need to make sure you also add its corresponding **Command spec**.

### Command spec

Command specs can be written in YAML (for now, but perhaps hcl, json, toml could come next). They should be named exactly like your script, minus the extension. For `.milpa/commands/my-command.sh` the corresponding spec file would be at `.milpa/commands/my-command.yaml`.

```yaml
# example spec
summary: Creates a github release # this shows up during autocomplete and command listings
description: | # here, you add a longer description of how your command does its magic
  This command shows you a changelog and waits for approval before generating and pushing a new tag, creating a github release, and opening the browser at the new release.

  ## Schemes

  You may choose between [semver](https://semver.org) and [calver](https://calver.org). Their composition is as follows:

  - **semver**: `{major}.{minor}.{patch}`, i.e. `3.14.1592`
  - **calver**: `{date}.{minor}.{micro}`, where date is derived from the `prefix` option; for example `16.18.339`.

# Arguments to commands are specified as a list
arguments:
  # each argument gets a name, your script will find these in the environment
  # so in this example, increment becomes $MILPA_ARG_INCREMENT
  # but it's also available as the first argument (i.e. `$1` in bash)
  - name: increment
    # arguments can be validated to be part of a set
    set:
      # sets can be static, like so:
      values: [micro, patch, minor, major]
      # dynamic sets can come from the lines printed by another subcommand
      # in this example, from `milpa scm versions`
      # from: { subcommand: scm versions }
    # arguments can have a default
    default: patch
    # or be required
    required: false

# options or flags are specified as a map
options:
  # this sets the --scheme option
  # it will be available to your script as the $MILPA_OPT_SCHEME environment variable
  scheme:
    description: Determines the format of the tags for this repo.
    # options can also be validated against a set of values
    set: { values: [semver, calver] }
    # and they may have defaults
    default: semver
    # or be required
    required: false
  prefix:
    description: |
      An optional prefix to prepend to release identifiers. If `calver` is chosen as `scheme`, you may specify a combination of `YY`, `YYYY`, `MM`, and `DD` to be replaced with the corresponding values of the local date. The default in that case is `YY`.
    default: ""
  ask:
    description: Prompt for confirmation before creating a github release
    # flags can be boolean. Since environment variables can only be strings,
    # these false values (or defaults) will be set as an empty string, while
    # true values will be set as "true"
    type: boolean
    default: false
```

## Milpa has commands itself

```sh
# install autocomplete scripts for your $SHELL
milpa itself shell install-autocomplete
# add this to your shell profile to source environment variables used by your commands
milpa itself shell init
# create new milpa commands on the local repo
milpa itself create
```

## Example

See [unRob/nidito](https://github.com/unRob/nidito/tree/master/.milpa).

## Internals

Milpa is built with:

- [bash](https://www.gnu.org/software/bash/)
- [spf13/cobra](https://cobra.dev)
