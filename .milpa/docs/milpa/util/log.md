---
description: Functions related to output
---

The `log` util contains shell functions related to output of informational messages. This utility is **loaded by default**.

If [`TERM`](https://linux.die.net/man/7/term) is not set, then it'll default to `xterm-color`.

## Functions

### `@milpa.fmt`

`@milpa.fmt CODE MESSAGE`

Returns `MESSAGE` formatted with `CODE` and does not print a new line. Code can be any of `bold`, `warning`, `error` and `inverted`.

### `@milpa.log`

`@milpa.log LEVEL [MESSAGE]`

Prints `MESSAGE` with a log prefix, unless `--silent` is specified or `LEVEL` is error. Level may be one of:

- `complete`: prefixes ✅ to your message and prints in bold,
- `success`: prefixes ✔ to your message,
- `error`: prints your message in red,
- `warning`: prints your message in yellow,
- `info`: just prints your message, and
- `debug`: if `--verbose` is specified, prints your message in gray.
