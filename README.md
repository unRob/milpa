# milpa

[`milpa`](https://milpa.dev) is a command-line tool to care for one's own garden of scripts. [Its name](https://en.wikipedia.org/wiki/Milpa) comes from an agricultural method that combines multiple crops in close proximity.

For a brief introductory tutorial, check out [`milpa help docs milpa quick-guide`](/.milpa/docs/milpa/quick-guide.md).

```sh
# install on mac and linux with:
curl -L https://milpa.dev/install.sh | bash -
```

You and your team write scripts and a little spec for each of them—use bash, or any other language—, and `milpa` provides:

- argument and option **completions** from static sources, other `milpa` scripts or even other programs;
- **nested sub-commands** as simple as the filesystem, organize them in folders and `milpa` does the rest;
- **parsing** and **validation** for arguments and options, writing little to no code; and
- **help** and **documentation** on the terminal and browser.

There's [a few reasons why](/.milpa/docs/milpa/use-case.md) you and your team might wanna use `milpa`, but in summary, it's goal is to provide all those nice features above while making it easier to follow the [Command Line Interface Guidelines](https://clig.dev).

`milpa` is licensed under the Apache License 2.0, and its code is available at [github.com/unRob/milpa](https://github.com/unRob/milpa).

## Concepts

`milpa` runs **commands** found in one or more **repos**:

- **Commands** are bash scripts or executables written in your language of choice, and their corresponding specs written in YAML. Get the full story with [`milpa help docs milpa command`](/.milpa/docs/milpa/command/index.md).
- Commands are organized in folders within one or more **repos**. Repos are just folders that contain a `.milpa` folder within. Check out [`milpa help docs milpa repo`](/.milpa/docs/milpa/repo/index.md) for more details about repos.

## Enough words, show me some code

```sh
# milpa is a program you run
milpa
# Add a command at .milpa/commands/hello.{sh,yaml} and milpa will gladly run it
milpa hello
# How do you create a milpa command? Check out:
milpa itself create --help
# Speaking of, milpa comes with fancy documentation
milpa help docs milpa
# You can also browse through an HTML version of it
milpa help docs --server
# Trouble with milpa? see what's making milpa sad
milpa itself doctor

# Install whole groups of commands written by strangers on the internet!
milpa itself repo install github.com/nidito/unRob
# and use them like so:
# milpa nidito dc list

# When typing gets annoying, install autocomplete scripts for your $SHELL
milpa itself install-autocomplete
```
