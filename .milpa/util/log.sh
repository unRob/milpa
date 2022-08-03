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

function @milpa.is_color_enabled() {
  if [[ -n "${NO_COLOR}" ]]; then
    return 1
  fi

  [[ "$COLOR" == "always" ]] || [[ -t 2 ]];
}

_FMT_BOLD=$'\e[1m'
_FMT_DIM=$'\e[2m'
_FMT_UNDERLINE=$'\e[4m'
_FMT_INVERT=$'\e[7m'
_FMT_RESET=$'\e[0m'
_FMT_FG_DEFAULT=$'\e[39m'
_FMT_FG_RED=$'\e[31m'
_FMT_FG_GREEN=$'\e[32m'
_FMT_FG_YELLOW=$'\e[33m'
_FMT_FG_GRAY=$'\e[37m'
_FMT_BG_DEFAULT=$'\e[49m'
_FMT_BG_RED=$'\e[41m'
_FMT_BG_GREEN=$'\e[42m'
_FMT_BG_YELLOW=$'\e[43m'
_FMT_BG_GRAY=$'\e[47m'


function @milpa.fmt() {
  if ! @milpa.is_color_enabled; then
    shift
    echo -e "$*"
    return
  fi

  local code;
  case $1 in
    bold) code="$_FMT_BOLD" ;;
    warning) code="$_FMT_FG_YELLOW" ;;
    error) code="$_FMT_FG_RED" ;;
    inverted) code="$_FMT_INVERT" ;;
    underlined) code="$_FMT_UNDERLINE" ;;
    *) @milpa.fail "unknown formatting directive: $1" ;;
  esac
  shift
  echo -e "${code}$*${_FMT_RESET}"
}

function _print_message () {
  local level command_name
  level=$1
  shift
  date=""
  if [[ "$MILPA_VERBOSE" != "" ]]; then
    date=" $(date -u +"%FT%H:%M:%S")"
  fi
  command_name=${MILPA_COMMAND_NAME:-milpa}

  [[ "$level" == "debug" ]] && [[ -z "${MILPA_VERBOSE+x}${DEBUG+x}" ]] && return
  >&2 echo "${_FMT_DIM}[${level}:${command_name// /:}${date}]${_FMT_RESET} $*"
}

function @milpa.log () {
  local prefix format level;
  level="info"
  case "$1" in
    complete) prefix="✅ "; format="$_FMT_BOLD" ;;
    success) prefix="✔ "; format="$_FMT_BOLD" ;;
    error) level="error"; format="$_FMT_FG_RED" ;;
    warn*) level="warning"; format="$_FMT_FG_YELLOW" ;;
    info) level="info" ;;
    debug) level="debug"; format="$_FMT_DIM" ;;
    *)
      @milpa.log warn "unknown log kind: $1"
      level="bad-milpa-log"
      format="$_FMT_FG_RED"
      set -- "" "$@"
      ;;
  esac
  if [[ "$MILPA_SILENT" == "true" ]] && [[ "$level" != "error" ]]; then
    return
  fi
  shift
  msg="$*"
  if [[ $# == 0 ]]; then
    msg=$(cat)
  fi

  if @milpa.is_color_enabled && [[ -n ${format:+x} ]]; then
    msg="${format}$msg${_FMT_RESET}"
  fi

  _print_message "$level" "${prefix}$msg"
}
