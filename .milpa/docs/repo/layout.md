---
related-docs: [commands, environment]
related-commands: ["itself create"]
---
# milpa repo layout

Repositories are folders that contain a `.milpa` folder within. Use the `MILPA_PATH` environment variable to tell `milpa` where to look for repos (see [`itself docs environment`](docs/milpa/environment#MILPA_PATH)). By default, `milpa` will prepend any folder named `.milpa` at the top-level of a git repository to the `MILPA_PATH`.

Repositories must contain a `commands` folder, with [commands](/docs/milpa/commands), and may also include `utils` to be used by command executables, and `docs`, to document anything related to your `milpa` repo, and/or `git` repo.

## Example

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

```

Then, `milpa` would allow you to run `milpa vault cloud-provider login` and `milpa onboard`, as well as `milpa vault db connect api --environment production --verbose`, or even `milpa help vault db list` and so on and so forth. You choose how to organize your milpa commands under `.milpa/commands`, and `milpa` figures out the rest.
