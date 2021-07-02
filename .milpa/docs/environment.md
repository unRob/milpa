---
related-docs: [utils/log, command/environment]
---

# milpa environment

There's a few environment variables that control the behavior of `milpa`.

## Paths

### MILPA_ROOT

`MILPA_ROOT` points to the installed milpa _kernel_, by default `/usr/local/lib/milpa`. This folder contains a milpa repo, the `milpa` executable, and a helper binary named `compa`, along a copy of the license and the source repo's README.

You can set it to a local installation, like a fork, and run `$MILPA_ROOT/milpa` to use that fork's scripts instead of an installed version.

If you choose to install to a different path than the default, you may wanna add `MILPA_ROOT=/path/to/your/install` you your shell profile.

### MILPA_PATH

The `MILPA_PATH` environment variable tells `milpa` where to look for repos. By default, `milpa` will prepend any folder named `.milpa` at the top-level of a git repository to the `MILPA_PATH`. If git is not available, it will look in your current working directory instead.

Additional repositories can be added as colon (`:`) delimited paths, pointing to the directory containing a `/.milpa` folder within. For example, setting `MILPA_PATH=$HOME/code/my-repo:$HOME/.local/milpa` would add the `$HOME/code/my-repo` and `$HOME/.local/milpa` folders to the command search path. Commands with the same name will override any commands previously found in the `MILPA_PATH`.

If desired, you may set a `MILPA_PATH` for all shells by adding it to your shell's profile.

## Output

### NO_COLOR

Also enabled by the `--no-color` option to disable printing of formatting escape codes from `compa` and `_log`.

### MILPA_VERBOSE

Enabled by the `--verbose` option. It shows information about what milpa is doing, along any `_log debug` messages from commands.

### MILPA_SILENT

Enabled by the `--silent` option, to hide `_log` messages completely

## DEBUG

Set `DEBUG=1` and find out whatever roberto needed to debug before writing proper tests.
