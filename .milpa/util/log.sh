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

if [[ -x ${TERM+x} ]]; then
  # otherwise, tput is gonna have trouble fetching formatting codes
  export TERM="xterm-color"
fi

if [[ -t 1 ]] && [[ -z ${NO_COLOR+x} ]]; then
  _FMT_INVERTED=$(tput rev)
  _FMT_BOLD="$(tput bold)"
  _FMT_RESET="$(tput sgr0)"
  _FMT_ERROR="$(tput setaf 1)"
  _FMT_WARNING="$(tput setaf 3)"
  _FMT_GRAY="$(tput setaf 7)"
else
  _FMT_INVERTED=""
  _FMT_BOLD=""
  _FMT_RESET=""
  _FMT_ERROR=""
  _FMT_WARNING=""
  _FMT_GRAY=""
fi


function @milpa.fmt() {
  local code;
  case $1 in
    bold) code="$_FMT_BOLD" ;;
    warning) code="$_FMT_WARNING" ;;
    error) code="$_FMT_ERROR" ;;
    inverted) code="$_FMT_INVERTED" ;;
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
  if [[ "$MILPA_VERBOSE" == 1 ]]; then
    date=" $(date -u +"%FT%H:%M:%S")"
  fi
  command_name=${MILPA_COMMAND_NAME:-milpa}

  [[ "$level" == "debug" ]] && [[ -z "${MILPA_VERBOSE+x}" ]] && return
  >&2 echo "${_C_GRAY}[${level}:${command_name// /:}${date}]${_FMT_RESET} $*"
}

function @milpa.log () {
  local prefix format level;
  level="info"
  case $1 in
    complete) prefix="✅ "; format="$_FMT_BOLD" ;;
    success) prefix="✔ "; format="$_FMT_BOLD" ;;
    error) level="error"; format="$_FMT_ERROR" ;;
    warn*) level="warning"; format="$_FMT_WARNING" ;;
    info) ;;
    debug) level="debug"; format="$_FMT_GRAY" ;;
    *) @milpa.fail "unknown log kind: $1" ;;
  esac
  if [[ "$MILPA_SILENT" == "true" ]] && [[ "$level" != "error" ]]; then
    return
  fi
  shift
  msg="$*"
  if [[ $# == 0 ]]; then
    msg=$(cat)
  fi

  if [[ -n ${format:+x} ]]; then
    msg="${format}$msg${_FMT_RESET}"
  fi

  _print_message "$level" "${prefix}$msg"
}
