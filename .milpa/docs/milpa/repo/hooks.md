---
title: milpa repo hooks
---
# !milpa! repo hooks

`!milpa!` provides a couple of hook points for you to tweak the behavior of your shell and your repo's commands. Hooks for your repo must be placed in the `.milpa/hooks` folder.

## `.milpa/hooks/shell-init(.sh)`

This hook is run whenever `!milpa! itself shell init` is called. Its purpose is to set any environment variables specific to your repo during every shell's initialization process. This hook can be either a bash shell script with an `.sh` extension, or an executable file without extension.

## `.milpa/hooks/before-run.sh`

This hook runs before invoking any command from your repo, and may be useful to do additional validations or checks before actually calling any of your commands. The full [environment](/.milpa/docs/milpa/environment.md) is available for this hook. This hook must be a bash shell script with an `.sh` extension.
