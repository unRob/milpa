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
