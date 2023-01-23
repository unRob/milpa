---
description: Functions related to temporary files
---

The `fmt` util contains shell functions related to the creation and cleanup of temporary files and directories.

## Functions

### `@tmp.file`

`@tmp.file HANDLE`

Creates a temporary file (in `/tmp` with `HANDLE` as prefix) and exports a new variable named exactly like the value of `HANDLE` pointing to the path of this new temporary file. `HANDLE` should therefore be a valid variable name identifier.

```sh
#!/usr/bin/env bash
@milpa.load_util tmp

# create the file
@tmp.file my_tmp_file

echo "some data" >"$my_tmp_file"
```

### `@tmp.dir`

`@tmp.dir PREFIX`

Creates a temporary directory (in `/tmp` with `HANDLE` as prefix) and exports a new variable named exactly like the value of `HANDLE` with the path to this new directory. `HANDLE` should therefore be a valid variable name identifier.


```sh
#!/usr/bin/env bash
@milpa.load_util tmp

# create the directory
@tmp.dir my_tmp_data

# write to it
echo "some data" >"$my_tmp_data/some.thing"
echo "other data" >"$my_tmp_data/something.else"
```

### `@tmp.cleanup`

`@tmp.cleanup`

Cleans up all temporary files created so far for all users of this utility.


```sh
#!/usr/bin/env bash
@milpa.load_util tmp
@tmp.file my_tmp_file
echo "some data" >"$my_tmp_file"

# delete all temporary files on normal exit or error
trap '@tmp.cleanup' ERR EXIT

# or run as needed
@tmp.cleanup
```

