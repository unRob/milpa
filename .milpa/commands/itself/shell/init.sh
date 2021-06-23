#!/usr/bin/env bash
_CURRENT_SHELL=$(basename "$SHELL")
function xsh_export (){
  case "$_CURRENT_SHELL" in
    fish)
      echo "set -x $1 $2"
      ;;
    *sh)
      echo "export $1=\"$2\""
      ;;
  esac
}

function xsh_prepend (){
  path_var="${2:-PATH}"
  case "$_CURRENT_SHELL" in
    fish)
      # set PATH some/new/path $PATH
      echo "set $path_var $1 \$$path_var"
      ;;
    *sh)
      # export PATH="some/new/path:$PATH"
      echo "export $path_var=\"$1:\$$path_var\""
      ;;
  esac
}

IFS=':' read -r -a args <<< "${MILPA_PATH//:/\/.milpa/hooks:}/.milpa/hooks"
_log debug "looking for env files in ${args[*]}"
args+=( -name "shell-init" -o -name "shell-init.sh" )
find "${args[@]}" 2>/dev/null | while read -r env_file; do
  if [[ -x "$env_file" ]]; then
    _log debug "executing $env_file"
    "$env_file"
  else
    _log debug "sourcing $env_file"
    # shellcheck disable=1090
    source "$env_file"
  fi
done
