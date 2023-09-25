#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>


function upgrade_repo() {
  @milpa.log info "Upgrading $(@milpa.fmt bold "$1")"

  cd "$1" || @milpa.fail "could not cd into $1"
  git pull --depth=1 origin "$(git branch --show-current)" || @milpa.fail "Could not upgrade $1"
  @milpa.log "${2:-complete}" "Upgraded $1"
}

if [[ "$MILPA_ARG_PATH" ]]; then
  upgrade_repo "$MILPA_ARG_PATH"
  exit
fi

@milpa.log info "Upgrading all cloned repos"
while read -r clone; do
  upgrade_repo "$clone" success
done < <(milpa itself repo list --cloned --paths-only)
@milpa.log complete "All cloned repos upgraded"
