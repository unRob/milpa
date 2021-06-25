#!/usr/bin/env bash
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
