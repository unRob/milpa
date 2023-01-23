---
description: Hook points for milpa commands
---

`milpa` provides a few hook points for you to tweak the behavior of your repo's commands. Hooks for your repo must be placed in the `.milpa/hooks` folder.


## `before-run.sh`

This hook runs before invoking any command from your repo, and may be useful to do additional validations or checks before actually calling any of your commands. The full [environment](/.milpa/docs/milpa/environment.md) is available for this hook. This hook must be a bash shell script with an `.sh` extension.

## `post-install.sh`

This hook is run after `milpa itself repo install` to bootstrap the installation of a remote repository (see [`milpa itself repo install`](/.milpa/commands/itself/repo/install.md)). This hook must be a bash shell script with an `.sh` extension.

## `post-uninstall.sh`

This hook is run during `milpa itself repo uninstall` to perform any necessary cleanup of a remote repository (see [`milpa itself repo uninstall`](/.milpa/commands/itself/repo/uninstall.md)). This hook must be a bash shell script with an `.sh` extension.
