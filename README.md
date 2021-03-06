# milpa

[`milpa`](https://milpa.dev) is a command-line tool to care for one's own garden of scripts, its name comes from ["milpa"](https://en.wikipedia.org/wiki/Milpa), an agricultural method that combines multiple crops in close proximity

For a brief intro to get started, check out [milpa help docs milpa quick-guide](/.milpa/docs/milpa/quick-guide.md).

```sh
# install on mac and linux with:
curl -L https://milpa.dev/install.sh | bash -
```

You and your team write scripts and a little spec for each of them—use bash, or any other language—, and `milpa` provides **autocompletions**, **subcommands**, **argument parsing** and **validation** so you can skip the toil and focus on your scripts.

There's [a few reasons why](/.milpa/docs/milpa/alternatives.md) you and your team might wanna use milpa, but basically, it's meant to provide all those nice features above while making it easier to follow the [Command Line Interface Guidelines](https://clig.dev).

`milpa` is licensed under the Apache License 2.0, and its code is available at [github.com/unRob/milpa](https://github.com/unRob/milpa).

## Concepts

**Commands** are bash scripts or executables written in your language of choice, and their corresponding specs. Get the full story with [milpa help docs milpa command](/.milpa/docs/milpa/command/index.md)

Commands are grouped into one or more **repos**. These are folders that contain a `.milpa` folder within. Check out [`milpa help docs milpa repo`](/.milpa/docs/milpa/repo/index.md) for more details about repos.

## Enough words, show me some code

```sh
# milpa is a program you run
milpa
# you can add a command at .milpa/commands/hello.{sh,yaml} and milpa will gladly run it
milpa hello
# when typing gets annoying, install autocomplete scripts for your $SHELL
milpa itself shell install-autocomplete
# see how to create new milpa commands on the local repo
milpa itself create --help
# read some documentation
milpa help docs milpa
# see what's making milpa sad
milpa itself doctor
# add new commands to milpa written by strangers on the internet!
milpa itself repo install github.com/nidito/unRob
```
