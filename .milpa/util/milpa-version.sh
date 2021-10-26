#!/usr/bin/env bash

MILPA_UPDATE_URL="${MILPA_UPDATE_URL:-https://milpa.dev/.well-known/milpa/latest-version}"
MILPA_UPDATE_PERIOD_DAYS="${MILPA_UPDATE_PERIOD_DAYS:-7}"
MILPA_UPDATE_PERIOD_SECONDS=$(( MILPA_UPDATE_PERIOD_DAYS * 24 * 3600 ))
MILPA_LOCAL_SHARE="${XDG_HOME_DATA:-$HOME/.local/share}/milpa"
_milpa_last_checked_path="${MILPA_LOCAL_SHARE}/last-update-check"


function @milpa.version.installed () {
  "$MILPA_COMPA" __version 2>&1 || {
    if [[ "$?" != 42 ]]; then
      @milpa.log debug "could not get installed version"
      return 1
    fi
  }
}

function @milpa.version.latest () {
  if ! curl --silent --fail -L --max-time 1 "$MILPA_UPDATE_URL"; then
    @milpa.log debug "Could not fetch latest version!"
    return 1
  fi
}

function @milpa.version.needs_check () {
  local now last_ping elapsed
  now=$(date +%s)
  last_ping="$(cat "$_milpa_last_checked_path" 2>/dev/null || echo 0)"
  elapsed=$(( now - last_ping ))

  @milpa.log debug "Looked for updates at $last_ping, $elapsed seconds ago"

  [[ "$elapsed" -ge "$MILPA_UPDATE_PERIOD_SECONDS" ]]
}

function @milpa.version.is_latest () {
  local installed latest

  if ! installed=$(@milpa.version.installed); then
    @milpa.log debug "Failed querying for current version"
    return 0
  fi

  if ! latest=$(@milpa.version.latest); then
    @milpa.log debug "Failed looking up latest version"
    return 0
  fi

  if [[ "$installed" == "$latest" ]]; then
    # keep bugging until this is true, i guess
    date "+%s" > "$_milpa_last_checked_path"
    return 0
  fi

  echo "$latest" "$installed"
  return 1
}

function _milpa_check_for_updates_automagically () {
  local versions latest installed
  [[ "$MILPA_COMMAND_NAME" == "itself upgrade" ]] && return 0

  if @milpa.version.needs_check && ! versions=$(@milpa.version.is_latest); then
    read -r latest installed <<<"$versions"
    MILPA_COMMAND_NAME=milpa @milpa.log warning "milpa $latest is available (you're running $installed), to upgrade run:"
    MILPA_COMMAND_NAME=milpa @milpa.log warning "milpa itself upgrade"
  fi
}
