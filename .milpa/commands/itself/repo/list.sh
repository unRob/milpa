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
    if [[ -f "$src/../.git/config" ]]; then
      if [[ $src = "${MILPA_REPOS_GLOBAL%%/repos}/clones"* ]] || [[ $src = "${MILPA_REPOS_USER%%/repos}/clones"* ]]; then
        src="git clone from $(cd "$src" && git remote get-url origin), repo at ${src%%/.milpa}"
      else
        src="local symlink, original at $src"
      fi
    fi
  elif [[ "$1" != "${MILPA_REPOS_USER}/"* ]] && [[ "$1" != "${MILPA_REPOS_GLOBAL}/"* ]]; then
    src="from \$MILPA_PATH"
  else
    src="source at local directory"
  fi
  echo "$(@milpa.fmt bold "$1") - $src"
}

if [[ "$MILPA_OPT_CLONED" ]]; then
  find -L "${MILPA_REPOS_USER%%/repos}/clones" "${MILPA_REPOS_GLOBAL%%/repos}/clones" -maxdepth 1 -mindepth 1 -type d 2>/dev/null | while read -r repo; do
    print_repo "$repo"
  done
  exit
fi

[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo "$(@milpa.fmt inverted " Local repos "): $MILPA_REPOS_USER"
find -L "$MILPA_REPOS_USER" -maxdepth 1 -mindepth 1 -type d | while read -r repo; do
  print_repo "$repo"
done


[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo
[[ ! "$MILPA_OPT_PATHS_ONLY" ]] && echo "$(@milpa.fmt inverted " Global repos "): $MILPA_REPOS_GLOBAL"
find -L "$MILPA_REPOS_GLOBAL" -maxdepth 1 -mindepth 1 -type d | while read -r repo; do
  print_repo "$repo"
done
