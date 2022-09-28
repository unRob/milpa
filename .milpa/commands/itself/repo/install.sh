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

if  [[ "$MILPA_ARG_TARGET" != "" ]]; then
  target="$MILPA_ARG_TARGET"
else
  user_target="${XDG_DATA_HOME:-$HOME/.local/share}/milpa"
  target="${user_target}/repos"
  if [[ "$MILPA_OPT_GLOBAL" ]]; then
    target="${MILPA_ROOT}/repos"
  else
    # sometimes milpa gets installed as root, but each user
    # should have it's own user repo folder
    mkdir -p "$target"
  fi
fi

base="${MILPA_ARG_SOURCE%%/.milpa/*}"

function symlink_local () {
  local base repo_name dst;
  base="$(cd "$1" && pwd)" || @milpa.fail "Could not "
  repo_name="${base##*/}"
  repo_name="${repo_name#.}"
  dst="$target/$repo_name"
  [[ -d "$dst" ]] && @milpa.fail "A repo named $repo_name already exists at $dst"
  ln -sfv "$base/.milpa" "$dst"
}

if [[ -d "$base/.milpa" ]]; then
  @milpa.log info "Local repository detected, symlinking..."
  symlink_local "$base" || @milpa.fail "Failed to symlink"
  @milpa.log success "Symlink created"
else
  @milpa.log info "Downloading repo..."
  new_repo=$("$MILPA_COMPA" __fetch "$base/.milpa" "$target") || @milpa.fail "Failed to download repo"
  @milpa.log success "Repo dowloaded"
  echo -n "$MILPA_ARG_SOURCE" > "$new_repo/downloaded-from"
fi


@milpa.log info "Running repo setup tasks"
if [[ -f "$new_repo/hooks/post-install.sh" ]]; then
  @milpa.log info "Running post-install hook"
  #shellcheck disable=1090,1091
  source "$new_repo/hooks/post-install.sh" || @milpa.fail "Could not run post-install hook to completion"
fi

@milpa.log complete "Repo installed at $new_repo"
