---
related-docs: [utils/log, command/environment]
description: An overview of all milpa environment variables
weight: 10
---
There's a few environment variables that control the behavior of `milpa`.

## Paths

### `MILPA_ROOT`

`MILPA_ROOT` points to the installed milpa _kernel_, by default `/usr/local/lib/milpa`. This folder contains a the built-in milpa repo, the `milpa` executable, and a copy of the license and the source repo's README.

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

`MILPA_DISABLE_GIT`, `MILPA_DISABLE_USER_REPOS` and `MILPA_DISABLE_GLOBAL_REPOS` each disable the corresponding command lookups when set to `true`.

---

## Output

### `DEBUG`

Set `DEBUG=on` to have `milpa` produce debugging output on its behavior. Set `DEBUG=trace` to see even more output.

### `MILPA_VERBOSE`

Enabled by the `--verbose` option. It shows information about what `milpa` is doing, along any `@milpa.log debug` messages from commands.

### `MILPA_SILENT`

Enabled by the `--silent` option, to hide `@milpa.log` messages completely. If `DEBUG` or `MILPA_VERBOSE` are enabled, these will override `MILPA_SILENT`.

### `NO_COLOR` / `COLOR`

By default, when stdout is a TTY, `milpa` will output color escape characters. If `NO_COLOR` is set to any non-empty string, then no color will be output. If not connected to a TTY, and `COLOR=always` is set, then color escape characters will be printed.

If either `--no-color` or `--color` options are provided, these will override both `COLOR` and `NO_COLOR` environment variables.

If `COLORTERM` is set to `truecolor`, 24-bit colors will be used instead of the default 256-color palette.

### `MILPA_HELP_STYLE`

`MILPA_HELP_STYLE` controls the theme to use when rendering help pages, and must be one of `auto`, `dark`, `light`, and `markdown`.

---

## Input

### `MILPA_SKIP_VALIDATION`

If enabled, validation will be skipped for arguments and options. Also enabled with `--skip-validation`. **Skipping validation may be unsafe**, but may be useful when validation depends on unavailable data or services.

---

## Auto-updates

### `MILPA_UPDATE_CHECK_DISABLED`

If set to a non-empty string, update checks will not be performed periodically before running commands.

### `MILPA_UPDATE_PERIOD_DAYS`

How often to check for new releases to `milpa` in days, by default `7`.

---

## Command Environment

Your commands will also have specific environment variables available, check out [milpa help docs milpa command](/.milpa/docs/milpa/command/index.md#arguments-options-and-environment).
