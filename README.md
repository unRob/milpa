# milpa

_milpa_ is an agricultural method that combines multiple crops in close proximity. `milpa` is a Bash script and tool to care for one's own garden of scripts.

You and your team write scripts and a little spec for each command. Use bash, or any other command, and `milpa` provides autocompletions, sub-commands and argument parsing+validation for you to focus on your scripts. For those in the know, it makes following these [Command Line Interface Guidelines](https://clig.dev/) easier.

> This repo has just been planted, and everything here is extremely experimental.

## Usage

`milpa [sub-command...] [options] [arguments]`


## Milpa has commands itself

```sh
# install autocomplete scripts for your $SHELL
milpa itself shell install-autocomplete
# add this to your shell profile to source environment variables used by your commands
milpa itself shell init
# create new milpa commands on the local repo
milpa itself create
# see what's making milpa sad
milpa itself doctor
```

## Example repo

See [unRob/nidito](https://github.com/unRob/nidito/tree/master/.milpa).

## Internals

Milpa is built with:

- [bash](https://www.gnu.org/software/bash/)
- [spf13/cobra](https://cobra.dev)
