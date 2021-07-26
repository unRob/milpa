---
description: Milpa comes with a shell library full of functions
weight: 20
---
`milpa` scripts that are written with bash may use any of the built-in or repo-specific shell utilities under the `.milpa/utils` folder.

## Usage

`@milpa.load_util UTIL_NAME...`

In your scripts, you may load utilities like so:

```sh
#!/usr/bin/env bash

# loads the repo util
@milpa.load_util repo

@milpa.log info "Current repo path is: $(@milpa.repo.current_path)"

# loads shell, as well as repo
@milpa.load_util repo shell
```

## Utilities

There's a few built-in utilities:

- [log](log): Loaded by default, has output-related functions
- [repo](repo): has functions related to milpa repositories
- [shell](shell): Has functions used by [`shell-init`](/.milpa/help/docs/milpa/repo/hooks/#shell-initsh) hook scripts.

## Add your own

To add your own to your repo, create a `util` folder and add `your-util-name.sh` under it. You may put in there whatever you want available for your shell scripts. Ideally, you'd put in there a few functions, whose name may start with `@milpa`, but are not required to be named in any particular way.

Then, you can:

```sh
@milpa.load_util your-util-name

# and use whatever functions you put in there
do_some_cool_stuff "$SHELL"
```
