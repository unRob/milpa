---
related-docs: [milpa/environment, milpa/command/spec]
related-commands: ["itself create"]
weight: 15
description: Commands overview
---
`milpa` can run two types of commands:

- bash scripts, with an `.sh` extension, or
- executables without an extension, written in whatever language you want.

## Spec

In order for `milpa` to recognize your commands, you'll need to make sure you also add its corresponding [command spec](/.milpa/docs/milpa/command/spec.md).

## How `milpa` invokes your command

`milpa` invokes your command with `source`, if it's a bash script with an `.sh` extension, and otherwise with `exec`. If your command does not have an extension, it must have the executable bit on (`chmod +x .milpa/commands/your-command`).

## Arguments and Options

The arguments and options passed in the command line will be parsed and validated according to your spec. If valid, your command will receive arguments as usual (`$1` and so on), without known options. Valid arguments and options will be available to your command as environment variables as defined below.

## Environment Variables

Along the [global environment variables](/.milpa/docs/milpa/environment.md), your command will have a the following environment variables available:

### `MILPA_COMMAND_*`

Your script has access to the following variables set by `milpa` after parsing arguments and running validations:

- `MILPA_COMMAND_NAME`: the space delimited name of your command, i.e. `db connect`;
- `MILPA_COMMAND_KIND`: either `source` for `.sh` scripts, or `exec` for executables;
- `MILPA_COMMAND_REPO`: the path to the repo containing this command, i.e. `/home/you/project`; and
- `MILPA_COMMAND_PATH`: the full path to the executable being called

### `MILPA_ARG_*`

Arguments specified on your spec will show up as environment variables with the `MILPA_ARG_` prefix, followed by the name set in your spec. Names will be all uppercase, and dashes will be turned into underscores. See the [command spec](/.milpa/docs/milpa/command/spec.md) for more information on arguments, and this abbreviated example below:

```yaml
# an example.yaml command spec like:
arguments:
  - name: greeting
  - name: full-name
    variadic: true
```

```sh
#!/bin/env bash
# With an example.sh script like:
title_case() {
  set ${*,,}
  echo ${*^}
}

echo "$MILPA_ARG_GREETING $(title_case "${MILPA_ARG_FULL_NAME[*]}")"
```

```sh
# when ran like this:
milpa example hello example world
# would output:
#> hello Example World
```

### `MILPA_OPT_*`

Options show up on the environment with the `MILPA_OPT_` prefix followed by the name in your spec. Names will be all uppercase, and dashes will be turned into underscores. **Boolean** type options have a special behavior, they'll be an empty string (`""`) if `false`, and `"true"` if `true`, so comparing them in bash is simpler (i.e. `if [[ "$MILPA_OPT_BOOL_FLAG" ]] `). See the [command spec](/.milpa/docs/milpa/command/spec.md) for more information on options, and this abbreviated example below:

```yaml
# let's add options to the example above
options:
  shout:
    type: boolean
    default: false
```

```sh
# we modify the example above
if [[ "$MILPA_OPT_SHOUT" ]]; then
  echo "$MILPA_ARG_GREETING ${MILPA_ARG_FULL_NAME[*]}!" | awk '{print toupper($0)}'
else
  echo "$MILPA_ARG_GREETING $(title_case "${MILPA_ARG_FULL_NAME[*]}")"
fi
```

```sh
# and run with
milpa example --shout hello loud boi
# to get
#> HELLO LOUD BOI!
```
