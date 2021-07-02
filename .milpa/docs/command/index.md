---
related-docs: [milpa/environment]
related-commands: ["itself create"]
---

# milpa commands

`milpa` can run two types of commands:

- bash scripts, with an `.sh` extension, or
- executables without an extension, in whatever language you want.

## Spec

In order for `milpa` to recognize your commands, you'll need to make sure you also add its corresponding [command spec](docs/milpa/command/spec).

## Your command itself

`milpa` invokes your command with `source`, if it's a bash script with an `.sh` extension, and otherwise with `exec`. If your command does not have an extension, it must have the executable bit on (`chmod +x .milpa/commands/your-command`).

The arguments and options passed by the user will be parsed and validated according to your spec, known options will be removed, and arguments will be passed to your command as typed. Unknown options will raise an error and your command will not be called.

## Environment Variables

Your command will have a the following environment variables available:

### MILPA_COMMAND_*

Your script has access to the following variables set by milpa after parsing arguments and running validations:

- `MILPA_COMMAND_NAME`: the space delimited name of your command, i.e. `db connect`;
- `MILPA_COMMAND_KIND`: either `source` for `.sh` scripts, or `exec` for executables;
- `MILPA_COMMAND_REPO`: the path to the repo containing this command, i.e. `/home/you/project`; and
- `MILPA_COMMAND_PATH`: the full path to the executable being called

### MILPA_ARG_*

Arguments specified on your spec will show up as environment variables with the `MILPA_ARG_` prefix, followed by the name set in your spec. Names will be all uppercase, and dashes will be turned into underscores.

### MILPA_OPT_*

Options show up on the environment with the `MILPA_OPT_` prefix followed by the name in your spec. Names will be all uppercase, and dashes will be turned into underscores. **Boolean** type options have a special behavior, they'll be an empty string (`""`) if `false`, and `"true"` if `true`, so comparing them in bash is simpler (i.e. `if [[ "$MILPA_OPT_BOOL_FLAG" ]] `).
