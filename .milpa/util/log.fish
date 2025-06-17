#!/usr/bin/env fish
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function @milpa.is_color_enabled
  if set -q NO_COLOR
    return 1
  end

  if test "$COLOR" = "always" -o -t 2
    return 0
  end
  return 1
end

set -g _FMT_BOLD '\e[1m'
set -g _FMT_DIM '\e[2m'
set -g _FMT_UNDERLINE '\e[4m'
set -g _FMT_INVERT '\e[7m'
set -g _FMT_RESET '\e[0m'
set -g _FMT_FG_DEFAULT '\e[39m'
set -g _FMT_FG_RED '\e[31m'
set -g _FMT_FG_GREEN '\e[32m'
set -g _FMT_FG_WHITE '\e[225m'
set -g _FMT_FG_YELLOW '\e[33m'
set -g _FMT_FG_GRAY '\e[37m'
set -g _FMT_BG_DEFAULT '\e[49m'
set -g _FMT_BG_RED '\e[41m'
set -g _FMT_BG_GREEN '\e[42m'
set -g _FMT_BG_YELLOW '\e[43m'
set -g _FMT_BG_GRAY '\e[47m'

function @milpa.fmt
  if ! @milpa.is_color_enabled
    echo -e "$argv[2..]"
    return
  end

  set -f code "";
  switch $argv[1]
  case bold
    set -f code "$_FMT_BOLD"
  case warning
    set -f code "$_FMT_FG_YELLOW"
  case error
    set -f code "$_FMT_FG_RED"
  case inverted
    set -f code "$_FMT_INVERT"
  case underlined
    set -f code "$_FMT_UNDERLINE"
  case *
    @milpa.fail "unknown formatting directive: $argv[1]"
  end

  echo -e "$code$argv[2..]$_FMT_RESET"
end

function _print_message
  set -f level $argv[1]
  set -f prefix ""
  set command_name "$MILPA_COMMAND_NAME"
  if test "$command_name" = ""
    set -f command_name "milpa"
  end

  if test -n "$MILPA_VERBOSE"
    set -f cmd_name $(string replace --all " " ":" "$command_name")
    set -f prefix "$_FMT_DIM$(date -u +"%FT%H:%M:%S") $level $cmd_name$_FMT_RESET "
  else if test "$level" = "error"
    set -f prefix "ERROR: "
    if @milpa.is_color_enabled
      set -f prefix "$_FMT_BG_RED$_FMT_FG_WHITE$_FMT_BOLD ERROR $_FMT_RESET "
    end
  end

  if test "$level" = "debug" -a \( -n "$MILPA_VERBOSE" -o -n" $DEBUG" \)
    return
  end
  echo -e "$prefix$argv[2..]" >&2
end

function @milpa.log
  set -f level ""
  set -f format ""
  switch "$argv[1]"
  case "complete"
    set -f prefix "✅ "
    set -f format "$_FMT_BOLD"
  case "success"
    set -f prefix "✔ "
    set -f format "$_FMT_BOLD"
  case "error"
    set -f level "error"
    set -f format "$_FMT_FG_RED"
  case "warn"
     set -f level "warning"
     set -f format "$_FMT_FG_YELLOW"
  case "info"
    set -f level "info"
  case "debug"
    set -f level "debug"
    set -f format "$_FMT_DIM"
  case "*"
      @milpa.log warn "unknown log kind: $1"
      set -f level "bad-milpa-log"
      set -f format "$_FMT_FG_RED"
      ;;
  end

  if test "$MILPA_SILENT" = "true" -a "$level" != "error"
    return
  end

  set -f msg "$argv[2..]"
  if test (count "$msg") -eq 0
    set msg $(cat)
  end

  if @milpa.is_color_enabled && test (string length "$format") -gt 0
    set -f msg "$format$msg$_FMT_RESET"
  end

  _print_message "$level" "$prefix$msg"
end
