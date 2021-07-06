---
title: Getting started
weight: 1
---
After milpa is installed, you'll find the best experience enabling command completion for your shell.

We'll assume you have a repo with your team's code at `~/src/my-teams-app`, but you can also test `!milpa!` out by creating any folder, wherever you like.

Once that's set up, you'd get started with `!milpa!` by:

1. installing milpa's autocompletion script to your shell, possibly restarting your session,
2. heading to your team's project folder
3. creating a new command, and editing it's spec
4. running the command to test it out

## installing autocomplete

```sh
# one of the best things about milpa is it's built-in support for shell completion
# this happens after you type "milpa " at the prompt and press the `tab` key
# milpa will offer suggestions on commands, arguments and options based on your
# commmands's spec. We'll install the required scripts by running:
milpa itself shell install-autocomplete

# After this is done, you'll likely need to reload your shell or open a new session/tab.
```

## Creating a new command

```sh
# Start by heading to your project's git repo
cd ~/src/my-teams-app

# let's say we wanna share a script that every new hire runs to setup their system.
# milpa helps you get started by running `milpa itself create` to create a new command
# Check out how to use it by running
milpa itself create --help

# Now we've seen the help page for the `itself create` command, let's go ahead and create our "developer-setup" command:
milpa itself create developer-setup

# A couple of new files will be created:
ls -lah .milpa/commands
# you'll see output like:
# -rw-r--r--  1 ada  staff  3141 Mar  14 15:09 developer-setup.sh
# -rw-r--r--  1 ada  staff   314 Mar  14 15:09 developer-setup.yaml
```

## Actually writing a command

```sh
# now, we'll begin by editing our spec
"$EDITOR" .milpa/commands/developer-setup.yaml
```

Our [spec](./milpa/docs/command/spec.md) should look like this:

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

Now, we can begin writing our [command](./milpa/docs/command/index.md):

```sh
#!/usr/bin/env bash
# .milpa/commands/developer-setup.sh

@milpa.log info "🚀 getting your computer ready for your ${MILPA_OPT_KIND:} exploits 🚀"

@milpa.log info "installing some packages..."
@milpa.log info "configuring credentials sources"
@milpa.log info "pleasing the compliance overlords, overladies, and overfolks"

@milpa log complete "Your system is ready to roll, $(whoami)!"
```

Once we're done, we can check milpa recognizes the command and can parse it's spec:

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

# if you don't see a checkmark next to your command's name, you'll get some
# feedback on what went wrong, check out `milpa help docs spec` for more info on
# how to write a command spec
```

## Running your new command

Once everything looks good, we should be ready to test our command out!

```sh
milpa developer-setup
# you'll see output like:
# [info:developer-setup] 🚀 getting your computer ready for your fullstack exploits 🚀
# ...
# [info:developer-setup] ✅ Your system is ready to roll, Ada!
```
