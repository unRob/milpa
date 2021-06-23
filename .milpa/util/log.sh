#!/usr/bin/env bash
if [[ -x ${TERM+x} ]]; then
  # otherwise, tput is gonna have trouble fetching formatting codes
  export TERM="xterm-color"
fi

_FMT_INVERTED=$(tput rev)
_FMT_BOLD="$(tput bold)"
_FMT_RESET="$(tput sgr0)"
_FMT_ERROR="$(tput setaf 1)"
_FMT_WARNING="$(tput setaf 3)"
_FMT_GRAY="$(tput setaf 7)"


function _fmt() {
  local code;
  case $1 in
    bold) code="$_FMT_BOLD" ;;
    warning) code="$_FMT_WARNING" ;;
    error) code="$_FMT_ERROR" ;;
    inverted) code="$_FMT_INVERTED" ;;
    *) _fail "unknown formatting directive: $1" ;;
  esac
  shift
  echo -e "${code}$*${_RESET}"
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
  >&2 echo "${_C_GRAY}[${level}:${command_name// /:}${date}]${_RESET} $*"
}

function _log () {
  local prefix format level;
  level="info"
  case $1 in
    complete) prefix="✅ "; format="$_FMT_BOLD" ;;
    success) prefix="✔ "; format="$_FMT_BOLD" ;;
    error) level="error"; format="$_FMT_ERROR" ;;
    warn*) level="warning"; format="$_FMT_WARNING" ;;
    info) ;;
    debug) level="debug"; format="$_FMT_GRAY" ;;
    *) _fail "unknown log kind: $1" ;;
  esac
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
