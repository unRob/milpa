---
description: Summary of how milpa works
weight: 50
---

`milpa` is a bash script in charge of setting environment variables and handing off to the user-requested script. Most of the heavy lifting, including argument/option parsing and validation, as well as finding commands, is done by a companion binary named `compa` (a slang term for friend in spanish).

`milpa` is built with, and thanks to:

- [bash](https://www.gnu.org/software/bash/)
- [spf13/cobra](https://cobra.dev)

## How it works

### `milpa` sets the stage

1. As it starts running, milpa will set `MILPA_ROOT` or exit unless it points to an existing directory,
2. We'll hand off to `compa` after setting the global environment. If the user requested a completion (with the hidden `__complete` command), the version flag, or the doctor/help docs server commands, `milpa` exits immediately after `exec compa $@`.
3. `@milpa.load_util` is defined otherwise, and we'll immediately use it to load logging-related functions.
4. `@milpa.fail` is defined.
5. A couple of temporary pipes are created (`COMPA_OUT` and `COMPA_ERR`), and compa is started


### `compa` resolves intentions

1. After setting up logging, `compa` builds and processes `MILPA_PATH` (unless `MILPA_PATH_PARSED` is already set), then looks for commands at `commands/` on every directory of `MILPA_PATH` and builds a command tree; it can error out here if any command has an invalid spec (unless running `milpa itself doctor`).
2. A `spf13/cobra.Command` is created and the known command tree is mapped into child commands.
3. `cobra` takes over, handling help, argument/flag parsing, and invoking validation.
4. Any errors are communicated back to milpa over the temporary pipes (`COMPA_OUT` and `COMPA_ERR`).
5. If a command is found, the user provided parseable arguments and options (and these are valid), the command environment is printed out to `COMPA_OUT`.

### `milpa` acts on the resolved intention

1. If a non-zero exit code is returned by `compa`, `milpa` will print out `compa`'s  `stdout`, piping it through `less -FIRX` if help or docs are being rendered, and exit with status code 0. Otherwise, we'll print out both pipes before removing them and exiting with `compa`'s original exit code.
2. If `compa` returned 0, `milpa` evals the contents of `COMPA_OUT` to set the found sub-command's environment. If an incomplete environment is found, we exit with status code 2, print debugging information and cleanup temporary pipes.
3. a version check is performed to nag the user to update to the latest available version
4. if requested, debug information of this session is printed to stderr, and temporary pipes are removed.
5. `before-run` hooks are ran before finally invoking your script.


## Exit codes

Mostly based on Bash's [Appendix E](https://tldp.org/LDP/abs/html/exitcodes.html) and FreeBSD's [`sysexits`](https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html)

| code  | reason |
|-------|--------|
| `2`   | `@milpa.fail` was called |
| `42`  | `compa` is requesting pretty printing and a clean `milpa` exit |
| `64`  | arguments/flags could not be parsed or failed validation |
| `70`  | a spec could not be parsed or help failed rendering |
| `78`  | `MILPA_ROOT` points to something that's not a directory, or `MILPA_PATH` has an incorrect path set |
| `127` | sub-command not found |

