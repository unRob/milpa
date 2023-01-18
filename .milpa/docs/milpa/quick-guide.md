---
title: Getting started
description: A quick guide to getting started with milpa
weight: 1
---

Getting started with `milpa` is a matter of following these steps:

0. Installing it,
1. installing milpa's autocomplete script to your shell, possibly restarting your session,
2. creating a new command, and editing it's spec
3. running the command to test it out

You'll be creating some files during this guide, so make sure to create a folder wherever you like. Feel free to use an existing project folder if you already have a script for it in mind!

## Steps

### Installing `milpa`

Let's begin by running the following command to install the latest version of `milpa`:

```sh
curl -L https://milpa.dev/install.sh | bash -
```

> If [homebrew](https://brew.sh) is available, you can also install `milpa` with:
> ```sh
>  brew install unRob/formulas/milpa
> ```

If all goes well, you'll see output that ends with something like:

```sh
-----------------------------------
🌽 Installed milpa version 3.1.4 🌽
-----------------------------------
Run 'milpa itself shell install-autocomplete' to install shell completions
```

You may run `milpa --version` to test `milpa` got installed correctly.

### Installing autocomplete

One of the best things about `milpa` is it's built-in support for shell completion. This is what happens after you type `milpa ` at a terminal prompt and press the `tab` key.

`milpa` will offer suggestions on commands, arguments and options based on your commands' spec. We'll install the required scripts by running:

```sh
milpa itself shell install-autocomplete
```

After this is done, you'll likely **need to reload your shell** or open a new session/tab. Verifying it's installed correctly can be done like so:

```sh
# type milpa, followed by a space, then press <TAB>
milpa
# you'll see this temporary output (on a zsh prompt, for this example):
# help      -- Display usage information on any command
# itself    -- itself subcommands
```

### Creating a new command

Let's assume we want to build a script that every new hire to our programming coop runs to setup their system. Once `milpa` is installed, we'll use `milpa itself create` to create a new [command](/.milpa/docs/milpa/command/index.md), composed of a bash script and a YAML [spec](/.milpa/docs/milpa/command/spec.md):

```sh
# Start by heading to your project's git repo
cd ~/src/my-teams-app

# milpa helps you get started by running `milpa itself create` to create a new command
# You can see how to use this command by running:
milpa itself create --help

# Now we've seen the help page for the `itself create` command,
# let's go ahead and create our "developer-setup" command:
milpa itself create developer-setup \
  --summary "Bootstrap a developer's machine" \
  --description "Installs all the fun stuff required for doing your dev job."
```

A couple of new files will be created, we verify with:

```sh
ls -lah .milpa/commands
# you'll see output like:
# -rw-r--r--  1 ada  staff  3141 Mar  14 15:09 developer-setup.sh
# -rw-r--r--  1 ada  staff   314 Mar  14 15:09 developer-setup.yaml
```

### Actually writing a command

Finally, we'll begin working on our command by editing our spec:

```sh
# open the yaml spec in your editor
"$EDITOR" .milpa/commands/developer-setup.yaml
```

Let's modify our spec at `.milpa/commands/developer-setup.yaml` to look like this:

```yaml
# .milpa/commands/developer-setup.yaml
# We'll set a one-line `summary` and a more elaborated `description`
summary: Bootstrap a developer's machine
description: |
  Installs all the fun stuff required for doing your dev job.

# We'll add a --kind option to specify which kind of work the user will be doing
options:
  kind:
    default: fullstack
    values: [fullstack, data]
    description: What kind of work will you be doing primarily?
```

Now, we should edit our script at `.milpa/commands/developer-setup.sh` to do fun stuff:

```sh
#!/usr/bin/env bash
# .milpa/commands/developer-setup.sh

@milpa.log info "🚀 getting your computer ready for your ${MILPA_OPT_KIND} exploits 🚀"

@milpa.log info "installing some packages..."
# then it's really up to you, here's some ideas:
# xcode-select --install
# brew bundle install --no-lock --file some.Brewfile
@milpa.log info "configuring credentials sources"
# aws configure ...
# gcloud auth ...
# op login ...
@milpa.log info "pleasing the compliance overlords, overladies, and overfolks"
# defaults write com.apple.SoftwareUpdate ScheduleFrequency -int 1
# defaults write com.apple.screensaver askForPassword -int 1
# defaults write com.apple.screensaver askForPasswordDelay -int 0

@milpa log complete "Your system is ready to roll, $(whoami)!"
```

Once we're happy with our developer setup script, we can check `milpa` recognizes the **command** and can parse it's **spec**:

```sh
milpa itself doctor
# output will look like:
# MILPA_ROOT is: /usr/local/lib/milpa
# MILPA_PATH is:
# /usr/local/lib/milpa
# /home/ada/src/my-teams-app
#
# Runnable commands:
# -----------
# ✅ developer-setup — /home/ada/my-teams-app/.milpa/commands/developer-setup.sh
# -----------
# ✅ itself create — /usr/local/lib/milpa/.milpa/commands/itself/create.sh
# -----------
# .. and so on
# -----------
```

If you don't see a checkmark next to your command's name, you'll get some feedback on what went wrong. Remember you can always run [`milpa help docs milpa command spec`](/.milpa/docs/milpa/command/spec.md) to see all the documentation on the command spec format.

### Running your new command

Once everything looks good, we should be ready to test our command out!

```sh
milpa developer-setup
# you'll see output like:
# [info:developer-setup] 🚀 getting your computer ready for your fullstack exploits 🚀
# ...
# [info:developer-setup] ✅ Your system is ready to roll, Ada!
```

Well done, `$(whoami)`! You've created and ran your first command. `milpa` comes with documentation on these commands and more, run [`milpa help docs milpa command`](/.milpa/docs/milpa/command/index.md) to find out more.

Finally, if you want to see _the rest of the owl_—that is, an actual implementation of this sort of setup script—then head over to [unRob/dotfiles](https://github.com/unRob/dotfiles/tree/master/.milpa/commands/computar) where you'll a practical example of how I setup my personal and work computers.

## Further reference

As you continue working on your scripts, make sure to check out these reference docs:

- [`milpa help docs milpa repo`](/.milpa/docs/milpa/repo/index.md): shows you how `milpa` expects you to organize your scripts
- [`milpa help docs milpa command`](/.milpa/docs/milpa/command/index.md): talks about the command structure and runtime
- [`milpa help docs milpa command spec`](/.milpa/docs/milpa/command/spec.md): detailed information on the command YAML spec
- [`milpa help docs milpa support`](/.milpa/docs/milpa/support.md): is the place where to ask for help or report a bug
- [`milpa help docs milpa environment`](/.milpa/docs/milpa/environment.md): context on the environment variables that modify `milpa`'s behavior.
