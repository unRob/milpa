#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function @milpa.load_util () {
  # shell scripts can call @milpa.load_util to load utils from MILPA_ROOT
  # or the current MILPA_COMMAND_REPO
  local env_name
  for util_name in "$@"; do
    env_name="_MILPA_UTIL_${util_name//-/_}"
    if [[ "${!env_name}" == "1" ]]; then
      # util already loaded
      continue
    fi

    global="$MILPA_ROOT/.milpa"
    util_path="${global}/util/$util_name.sh"
    if [[ ! -f "$util_path" ]] && [[ "$MILPA_COMMAND_REPO" != "" ]]; then
      util_path="$MILPA_COMMAND_REPO/util/$util_name.sh"
    fi

    if ! [[ -f "$util_path" ]]; then
      # util not found
      >&2 echo "Missing util named $util_name in ${MILPA_COMMAND_REPO}"
      exit 70 # programmer error
    fi

    set -o allexport
    # shellcheck disable=1090
    source "$util_path"
    set +o allexport
    export "${env_name?}=1"
    break
  done
}

@milpa.load_util log
function @milpa.fail () {
  # print an error, then exit
  @milpa.log error "$*"
  exit 2
}
