#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

MILPA_REPOS_USER="${XDG_DATA_HOME:-$HOME}/.local/share/milpa/repos"
MILPA_REPOS_GLOBAL="${MILPA_ROOT}/repos"

function print_repo() {
  local src
  if [[ "$MILPA_OPT_PATHS_ONLY" ]]; then
    echo "$1"
    return
  fi

  if [[ -L "$1" ]]; then
    src="$(readlink "$1")"
  elif [[ "$1" != "${MILPA_REPOS_USER}/"* ]] && [[ "$1" != "${MILPA_REPOS_GLOBAL}/"* ]]; then
    src="from \$MILPA_PATH"
  else
    src="$(cat "$1/downloaded-from")"
  fi
  echo "$(@milpa.fmt bold "$1") - $src"
}

[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo "$(@milpa.fmt inverted " Local repos "): $MILPA_REPOS_USER"
find -L "$MILPA_REPOS_USER" -maxdepth 1 -mindepth 1 -type d | while read -r repo; do
  print_repo "$repo"
done


[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo
[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo "$(@milpa.fmt inverted " Global repos "): $MILPA_REPOS_GLOBAL"
find -L "$MILPA_REPOS_GLOBAL" -maxdepth 1 -mindepth 1 -type d | while read -r repo; do
  print_repo "$repo"
done
