#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

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
