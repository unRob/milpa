#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

MILPA_REPOS_USER="${XDG_DATA_HOME:-$HOME}/.local/share/milpa/clones"
MILPA_REPOS_GLOBAL="${MILPA_ROOT}/repos"

@milpa.log info "Removing $(@milpa.fmt bold "$MILPA_ARG_PATH")"
if [[ -L "$MILPA_ARG_PATH" ]]; then
  src="$(dirname "$(readlink "$MILPA_ARG_PATH")")"
  if [[ $src = "${MILPA_REPOS_GLOBAL}/"* ]] || [[ $src = "${MILPA_REPOS_USER}/"* ]]; then
    @milpa.log info "removing cloned source at $src"
    rm -rf "$src"
  fi
  rm -f "$MILPA_ARG_PATH"
else
  rm -rf "$MILPA_ARG_PATH"
fi

if [[ -f "$MILPA_ARG_PATH/hooks/post-uninstall.sh" ]]; then
  @milpa.log info "Running post-uninstall hook"
  # run in a subshell so we don't care if it uninstall hook does weird stuff
  (
    #shellcheck disable=1090,1091
    source "$MILPA_ARG_PATH/hooks/post-uninstall.sh"
  ) || @milpa.log warning "Could not run post-uninstall hook to completion"
fi

@milpa.log complete "$MILPA_ARG_PATH uninstalled"
