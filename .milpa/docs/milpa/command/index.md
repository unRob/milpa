---
related-docs: [milpa/environment, milpa/command/spec]
related-commands: ["itself create"]
weight: 15
description: Commands overview
---

`milpa` runs your scripts or executables, by matching passed arguments to files in `.milpa/commands` directories. It looks at a few of these, starting with the one where `milpa` was called from. Check out `MILPA_PATH` of [`milpa help docs milpa environment`](/.milpa/docs/environment.md#MILPA_PATH) to understand where `milpa` looks for **commands**. A collection of commands under a single `.milpa` folder is called a **repo**, and you can learn more reading [`milpa help docs milpa repo`](/.milpa/docs/milpa/repo/index.md)

Your **script** or executable plus its corresponding **spec** is what we call a `milpa` **command**. Scripts can be:

0. bash scripts, with an `.sh` extension, or
1. executable files without an extension, written in whatever language you want. If your command does not have an extension, remember to set on the executable bit (`chmod +x .milpa/commands/your-command`)!

## Spec

A **spec** is a file named just like your **script**, except with a `.yaml` extension, with a few basic details about what your command does and what kinds of input it deals with. There's more details at [`milpa help docs milpa command spec`](/.milpa/docs/milpa/command/spec.md), but here's a brief example:

```yaml
# .milpa/commands/greet.yaml
summary: a very well mannered program
description: greets you politely on stdout
arguments:
  - name: full-name
    description: the name to greet
    variadic: true
options:
  greeting:
    default: Quihúbole mi
    description: the greeting word to use
```

## Script

A **script** then, will be able to use the definitions from your **spec** and do fun stuff with them, for example:

```sh
#!/bin/env bash
# .milpa/commands/greet.sh
function title_case() {
  set ${*,,}
  echo ${*^}
}

# arguments are passed as environment variables
if [[ "${#MILPA_ARG_FULL_NAME}" -eq  0 ]]; then
  name=$(whoami)
else
  name="${MILPA_ARG_FULL_NAME[*]}"
fi

# and so are options!
echo "$MILPA_OPT_GREETING $(title_case "$name")"
```

## Arguments, Options and Environment

The **arguments** and **options** passed in the command line will be parsed and validated according to your spec. If valid, your script will receive numbered arguments as usual (`$1` and so on), excluding any supplied options (i.e. without `--option=value` "arguments"). Valid arguments and options will be available to your script as **environment variables** as well.

> ⚠️ Unknown arguments and options will cause an error before your script executes!

The environment available to a **script** is composed of four groups:

0. `MILPA_COMMAND_*` variables have information about the **command** called by the user,
1. `MILPA_ARG_*` variables hold values for every **argument** of the spec,
2. `MILPA_OPT_*` variables hold values for every **option** defined, and
3. [global environment variables](/.milpa/docs/milpa/environment.md) that affect `milpa`'s overall behavior.


### Command metadata: `MILPA_COMMAND_*`

Your script has access to the following variables set by `milpa` after parsing arguments and running validations:

- `MILPA_COMMAND_NAME`: the space delimited name of your command, i.e. `db connect`;
- `MILPA_COMMAND_KIND`: either `source` for `.sh` scripts, or `exec` for executables;
- `MILPA_COMMAND_REPO`: the path to the repo containing this command, i.e. `/home/you/project/.milpa`; and
- `MILPA_COMMAND_PATH`: the full path to the executable being called.

### Arguments: `MILPA_ARG_*`

**Arguments** specified on your spec will show up as environment variables with the `MILPA_ARG_` prefix, followed by the name set in your spec. Names will be all uppercase, and dashes will be turned into underscores. See the [command spec](/.milpa/docs/milpa/command/spec.md) for more information on **arguments**, and this abbreviated example below:


```sh
# when ran like this:
milpa greet
# would output:
#> Quihúbole mi Grace

# you can also pass a specific name to greet, overriding the default
milpa greet elmer homero
# would output:
#> Quihúbole mi Elmer Homero
```

### Options: `MILPA_OPT_*`

**Options** show up on the environment with the `MILPA_OPT_` prefix followed by the name in your spec. Names will be all uppercase, and dashes will be turned into underscores. In the [command spec](/.milpa/docs/milpa/command/spec.md) you'll find more details on how to define **options**.

```sh
# when ran like this:
milpa greet --greeting Sup
# would output:
#> Sup Grace

# you can also pass a specific name to greet, overriding the default
milpa greet --greeting Oi joão
# would output:
#> Oi João
```

**Boolean** type options have a special behavior, their value will be an empty string (`""`) if `false`, and `"true"` if `true`, so comparing them in bash is simpler (i.e. `if [[ "$MILPA_OPT_BOOL_FLAG" ]] `), for example:

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
  echo "$MILPA_OPT_GREETING $(title_case "$name")!" | awk '{print toupper($0)}'
else
  echo "$MILPA_OPT_GREETING $(title_case "$name")"
fi
```

```sh
# and run with
milpa greet --shout loud boi
# to get
#> QUIHUBO MI LOUD BOI!
```
