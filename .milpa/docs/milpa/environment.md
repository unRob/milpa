---
related-docs: [utils/log, command/environment]
description: An overview of all milpa environment variables
weight: 10
---
There's a few environment variables that control the behavior of `milpa`.

## Paths

### `MILPA_ROOT`

`MILPA_ROOT` points to the installed milpa _kernel_, by default `/usr/local/lib/milpa`. This folder contains a milpa repo, the `milpa` executable, and a helper binary named `compa`, along a copy of the license and the source repo's README.

You can set it to a local installation, like a fork, and run `$MILPA_ROOT/milpa` to use that fork's scripts instead of an installed version.

If you choose to install to a different path than the default, you may wanna add `MILPA_ROOT=/path/to/your/install` you your shell profile.

### `MILPA_PATH`

The `MILPA_PATH` environment variable tells `milpa` where to look for repos.

By default, `milpa` will look for repos in the following order:

1. If `MILPA_PATH` is present in the environment, it'll start its search there,
2. then, `milpa` will look at its own commands under `$MILPA_ROOT`,
3. If the current working directory (or git repository) contains a .milpa folder, `milpa` will search that next,
4. from there, it'll look for user repos at `$XDG_DATA_HOME/milpa/repos`,
5. followed by global repos at `$MILPA_ROOT/repos`.

Commands with the same name will be silently ignored.

Additional repositories can be added as colon (`:`) delimited paths, pointing to the directory containing a `/.milpa` folder within. For example, setting `MILPA_PATH=$HOME/code/my-repo:/opt/milpa` would prepend the `$HOME/code/my-repo` and `/opt/milpa` folders to the command search path.

If desired, you may set a `MILPA_PATH` for all shells by adding it to your shell's profile.

---

## Output

### `DEBUG`

Set `DEBUG=1` and find out whatever roberto needed to debug before writing proper tests.

### `MILPA_VERBOSE`

Enabled by the `--verbose` option. It shows information about what `milpa` is doing, along any `@milpa.log debug` messages from commands.

### `MILPA_SILENT`

Enabled by the `--silent` option, to hide `@milpa.log` messages completely

### `NO_COLOR`

Also enabled by the `--no-color` option to disable printing of formatting escape codes from `compa` and `@milpa.log`.

---

## Command Environment

Your commands will also have specific environment variables available, check out [milpa help docs milpa command](/.milpa/docs/milpa/command/index.md#environment-variables)
