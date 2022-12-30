#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function @milpa.ask () {
  local prompt default result
  prompt="$1"
  if [[ "$2" ]]; then
    default="$2"
    prompt="$prompt [default: $default]"
  fi
  read -re -p "$prompt " result

  if [[ "$result" ]] || [[ "$default" ]]; then
    echo "${result:-$default}"
  else
    @milpa.warning "No value was entered, please try again."
    @ask "$prompt" "$result" "$default"
  fi
}

function @milpa.confirm () {
  read -r -p "$1${1:+ }Enter 'y' to continue: " -n 1
  [[ $REPLY =~ ^[Yy]$ ]]
  ret="$?"
  >&2 echo
  return "$ret"
}

function @milpa.select () {
  local options
  IFS=$'\n' read -r -d '' -a options <<<"$1"
  option_count=${#options[*]}

  PS3="Select an option (1-$(( option_count+1 ))): "
  select opt in "${options[@]}" "Quit"; do
    if [[ "$opt" == "Quit" ]] || [[ $REPLY == "$(( option_count + 1 ))" ]]; then
      return 1
    fi

    if [[ "$REPLY" != "" ]] && [[ "$REPLY" -gt 0 ]] && [[ "$REPLY" -le "$option_count" ]]; then
      echo "${opt}"
      break
    fi
    >&2 echo "No such option, try again"
  done
}
