#!/usr/bin/env fish
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function @milpa.load_util
  for util_name in $argv
    set env_name "_MILPA_UTIL_$(string replace --all "-" "_" "$util_name")"
    if set -q "$env_name"
      echo "util $util_name already loaded" >&2
      continue
    end

    set --local util_path "$MILPA_ROOT/.milpa/util/$util_name.fish"
    if test ! -f "$util_path" -a "$MILPA_COMMAND_REPO" != ""
      set util_path "$MILPA_COMMAND_REPO/util/$util_name.fish"
    end

    if ! test -f "$util_path"
      echo "Missing util named $util_name, add to $util_path" >&2
      exit 70
    end

    source "$util_path"
    set -g "$env_name" 1
  end
end

@milpa.load_util log

function @milpa.fail
  # print an error, then exit
  @milpa.log error $argv
  exit 2
end
