---
title: milpa docs
---
# milpa

"milpa", is an agricultural method that combines multiple crops in close proximity. [`milpa`](https://milpa.dev) is a Bash script and tool to care for one's own garden of scripts.

You and your team write scripts and a little spec for each command. Use bash, or any other command, and `milpa` provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.

For those in the know, it makes following these [Command Line Interface Guidelines](https://clig.dev/) easier.

## Concepts

### Repos

`milpa` let's you have one or more "repos", or collections of scripts. These are folders that contain a `.milpa` folder within. Check out [`milpa itself docs repo`](/.milpa/docs/milpa/repo/index.md) for more details about repositories.

### Commands

Your milpa repos will contain commands and their corresponding specs, among other things. These can be executables written in whatever language you want, or bash shell scripts. Get the full story with [milpa itself docs command](/.milpa/docs/milpa/command/index.md)

### Milpa comes with commands itself

```sh
# install autocomplete scripts for your $SHELL
milpa itself shell install-autocomplete
# see how to create new milpa commands on the local repo
milpa itself create --help
# read some documentation
milpa help docs milpa
# see what's making milpa sad
milpa itself doctor
```

## Enough words, show me some code

For a brief intro to get started, check out [milpa itself docs quick-guide](/.milpa/docs/milpa/quick-guide.md).

See [unRob/nidito](https://github.com/unRob/nidito/tree/master/.milpa) for a working example, and/or the drop that that spilled the cup and got me to work on `milpa`.

## Internals

Milpa is built with, and thanks to:

- [bash](https://www.gnu.org/software/bash/)
- [spf13/cobra](https://cobra.dev)

