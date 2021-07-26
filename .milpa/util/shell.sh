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
function @milpa.shell.export (){
  case "$_CURRENT_SHELL" in
    fish)
      echo "set -x $1 $2"
      ;;
    *sh)
      echo "export $1=\"$2\""
      ;;
  esac
}

function @milpa.shell.append_path (){
  path_var="${2:-PATH}"
  case "$_CURRENT_SHELL" in
    fish)
      # set PATH $PATH some/new/path
      echo "set $path_var \$$path_var $1"
      ;;
    *sh)
      # export PATH="$PATH:some/new/path"
      echo "export $path_var=\"\$$path_var:$1\""
      ;;
  esac
}

function @milpa.shell.prepend_path (){
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
