---
description: Milpa comes with a shell library full of functions
weight: 20
---

`milpa` scripts that are written with bash may use any of the built-in utilities, or any utilities defined in the same repo, under the `.milpa/utils` folder.

These are also bash scripts, with an `.sh` extension that provide functions that may be used by more than one command.


## Usage

`@milpa.load_util UTIL_NAME...`

In your scripts, you may load utilities like so:

```sh
#!/usr/bin/env bash

# loads the built-in util named `repo`
@milpa.load_util repo

@milpa.log info "Current repo path is: $(@milpa.repo.current_path)"

# load multiple utils at a time
@milpa.load_util repo shell
```

## Built-in utilities

There's a few utilities that come built-in to `milpa` and may be used by any bash scripts:

- [log](log): Loaded by default, has output-related functions
- [repo](repo): has functions related to milpa repositories
- [shell](shell): has functions used by [`shell-init`](/.milpa/docs/milpa/repo/hooks/#shell-initsh) hook scripts
- [tmp](tmp): has functions related to temporary files

## Add your own

To add your own utilities to your repo, create a `util` folder in it (that is`.milpa/util`) and add `your-util-name.sh` under it (`.milpa/util/your-util-name.sh`). You may put in there whatever you want available for your shell scripts.

### Example

Let's create a util at `.milpa/utils/voice.sh`:

```sh
# .milpa/utils/voice.sh
function yell() {
  echo "$@" | awk '{print toupper($0)}'
}

function yell-audibly() {
  osascript -e "say $(yell "$@")"
}
```

Then, on any of your repo's scripts, you can:

```sh
#!/usr/bin/env bash
@milpa.load_util voice

# and use whatever functions you put in there
yell "i <3 $SHELL"
```
