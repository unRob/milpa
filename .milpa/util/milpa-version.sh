#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

# The @milpa.version internal utils are abstract interactions with milpa's versions

# We check for milpa updates from this URL that returns the last-known version
# undocumented since this is only useful for tests and forks. If you fork milpa
# then might as well set this upstream or change this file (and bootstrap.sh).
MILPA_UPDATE_URL="${MILPA_UPDATE_URL:-https://milpa.dev/.well-known/milpa/latest-version}"
# check for new versions at most once every MILPA_UPDATE_PERIOD_DAYS days
MILPA_UPDATE_PERIOD_DAYS="${MILPA_UPDATE_PERIOD_DAYS:-7}"
MILPA_UPDATE_PERIOD_SECONDS=$(( MILPA_UPDATE_PERIOD_DAYS * 24 * 3600 ))

MILPA_LOCAL_SHARE="${XDG_HOME_DATA:-$HOME/.local/share}/milpa"
_milpa_last_checked_path="${MILPA_LOCAL_SHARE}/last-update-check"
# sometimes milpa gets installed as root, but each user should get its own dir
# to store update checkpoints
[[ -d "${MILPA_LOCAL_SHARE}" ]] || mkdir -p "$MILPA_LOCAL_SHARE"

function @milpa.version.log () {
  MILPA_COMMAND_NAME=milpa @milpa.log "$@"
}

# prints out the installed version
function @milpa.version.installed () {
  DEBUG=0 "$MILPA" --version 2>&1 || {
    if [[ "$?" != 42 ]]; then
      @milpa.version.log debug "could not get the installed version"
      return 1
    fi
  }
}

# prints out the latest known version
function @milpa.version.latest () {
  # since this is called before running milpa we need to timeout and just
  # keep going if version check takes longer than a second, but keep it configurable
  # so `milpa itself upgrade` has less of a chance to say there's no upgrade
  # available over slow connections. True story.
  local timeout; timeout="${1:-1}"

  if ! curl --silent --show-error --fail -L --max-time "$timeout" "$MILPA_UPDATE_URL"; then
    @milpa.log debug "Could not fetch latest version!"
    return 1
  fi
}

# tells if it's time for a new version check
function @milpa.version.needs_check () {
  local now last_ping elapsed
  now=$(date +%s)
  last_ping="$(cat "$_milpa_last_checked_path" 2>/dev/null || echo 0)"
  elapsed=$(( now - last_ping ))

  @milpa.version.log debug "Looked for updates at $last_ping, $elapsed seconds ago"

  [[ "$elapsed" -ge "$MILPA_UPDATE_PERIOD_SECONDS" ]]
}

# tells if the latest version is installed
# outputs the current and latest versions to stdout
function @milpa.version.is_latest () {
  local installed latest

  if ! installed=$(@milpa.version.installed); then
    @milpa.version.log debug "Failed querying for current version"
    return 0
  fi

  if ! latest=$(@milpa.version.latest "${1:-1}"); then
    @milpa.version.log debug "Failed looking up latest version"
    return 0
  fi

  # keep bugging until either the installed version is later or equal to latest
  if [[ "$installed" == "$latest" ]] || [[ "$installed" > "$latest" ]]; then
    date "+%s" > "$_milpa_last_checked_path"
    return 0
  fi

  echo "$latest" "$installed"
  return 1
}

# called by `milpa` to check for updates on run
function _milpa_check_for_updates_automagically () {
  local versions latest installed
  if [[ "$MILPA_COMMAND_NAME" == "itself upgrade" ]] || [[ "${MILPA_UPDATE_CHECK_DISABLED}" != "" ]]; then
    return 0
  fi

  if @milpa.version.needs_check && ! versions=$(@milpa.version.is_latest 1); then
    read -r latest installed <<<"$versions"
    @milpa.version.log warning "milpa $latest is available (you're running $installed), to upgrade run:"
    @milpa.version.log warning "milpa itself upgrade"
  fi
}
