---
description: Functions related to shells
---

The `shell` util contains shell functions useful for programs that output commands to be evaluated by a user's shell.

These functions are compatible with shells whose name ends in "sh" and are POSIX-compliant, and, the `fish` shell.

## Functions

### `@milpa.shell.export`

`@milpa.shell.export NAME VALUE`

Prints a command, that when evaluated by a user's shell, will set an environment variable on the current process, much like POSIX's `export` builtin.

### `@milpa.shell.append_path`

`@milpa.shell.export DIRECTORY [VARIABLE]`

Prints a command, that when evaluated by a user's shell, will append the `DIRECTORY` to `VARIABLE` (by default, `PATH`).

### `@milpa.shell.prepend_path`

`@milpa.shell.export DIRECTORY [VARIABLE]`

Prints a command, that when evaluated by a user's shell, will prepend the `DIRECTORY` to `VARIABLE` (by default, `PATH`).
