---
title: milpa repos
related-docs: ["milpa commands", "milpa environment"]
related-commands: ["itself create"]
weight: 10
---
Repositories are folders that contain a `.milpa` folder within. Use the `MILPA_PATH` environment variable to tell `!milpa!` where to look for repos (see [`!milpa! itself docs environment`](/.milpa/docs/milpa/environment.md#MILPA_PATH)). By default, `!milpa!` will prepend any folder named `.milpa` at the top-level of a git repository to the `MILPA_PATH`.

Repositories must contain a `commands` folder, with [commands](/.milpa/docs/milpa/command/index.md), and may also include `utils` to be used by command executables, [hooks](/.milpa/docs/milpa/repo/hooks.md) that modify the environment of `!milpa!` commands, and [docs](/.milpa/docs/milpa/repo/docs.md), to document anything related to your `!milpa!` repo.

Finally, milpa provides commands under [`milpa itself repo`](/.milpa/commands/itself/repo/index.md) to manage install, list and uninstall repositories from remote sources.

## Example repository layout

Let's say you have this in your repo:

```yaml
.milpa/
  commands/
    vault/
      cloud-provider/
        login
        login.yaml
      db/
        connect.sh
        connect.yaml
        list.sh
        list.yaml
    vpn/
      connect
      connect.yaml
    onboard.sh
    onboard.yaml
    release.sh
    release.yaml
  docs/
    welcome.md
    sdlc/
      releasing.md
  hooks/
    before-run.sh
    shell-init.sh
  util/
    date.sh
    etc.sh
    github.sh
```

Then, `!milpa!` would allow you to run `!milpa! vault cloud-provider login` and `!milpa! onboard`, as well as `!milpa! vault db connect api --environment production --verbose`, or even `!milpa! help vault db list` and so on and so forth. You choose how to organize your !milpa! commands under `.milpa/commands`, and `!milpa!` figures out the rest. Reading your welcome docs is as easy as `!milpa! help docs welcome`.

Using `@milpa.load_util` your posix-compliant shell scripts will be able to use any utils anywhere in the `MILPA_PATH`, for example, you could `@milpa.load_util github` and use any github-related functions in any of your repo's milpa commands.

Before any command runs, `.milpa/hooks/before-run.sh` would be called, and `.milpa/hooks/shell-init.sh` would be ran by `!milpa! itself shell init` during your shell's initialization process. See [hooks](/.milpa/docs/milpa/repo/hooks.md).

Ideally, you'll only store milpa-related files in your `.milpa` repo, as adding more files (specifically to the `commands` folder, will impact performance).
