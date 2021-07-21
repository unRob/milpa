`!milpa!` provides a couple of hook points for you to tweak the behavior of your shell and your repo's commands. Hooks for your repo must be placed in the `.milpa/hooks` folder.

## `shell-init(.sh)`

This hook is run whenever `!milpa! itself shell init` is called. Its purpose is to set any environment variables specific to your repo during every shell's initialization process. This hook can be either a bash shell script with an `.sh` extension, or an executable file without extension.

## `before-run.sh`

This hook runs before invoking any command from your repo, and may be useful to do additional validations or checks before actually calling any of your commands. The full [environment](/.milpa/docs/milpa/environment.md) is available for this hook. This hook must be a bash shell script with an `.sh` extension.


## `post-install.sh`

This hook is run after `!milpa! itself repo install` to bootstrap the installation of a remote repository (see [`milpa itself repo install`](/.milpa/commands/itself/repo/install.md)). This hook must be a bash shell script with an `.sh` extension.

## `post-uninstall.sh`

This hook is run during `!milpa! itself repo uninstall` to perform any necessary cleanup of a remote repository (see [`milpa itself repo uninstall`](/.milpa/commands/itself/repo/uninstall.md)). This hook must be a bash shell script with an `.sh` extension.
