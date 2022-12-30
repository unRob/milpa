#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

@milpa.load_util shell

IFS=':' read -r -a args <<< "${MILPA_PATH//:/\/hooks:}/hooks"
@milpa.log debug "looking for env files in ${args[*]}"
args+=( -name "shell-init" -o -name "shell-init.sh" )
find "${args[@]}" 2>/dev/null | while read -r env_file; do
  if [[ -x "$env_file" ]]; then
    @milpa.log debug "executing $env_file"
    "$env_file"
  else
    @milpa.log debug "sourcing $env_file"
    # shellcheck disable=1090
    source "$env_file"
  fi
done
