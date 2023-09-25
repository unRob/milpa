#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

user_target="${XDG_DATA_HOME:-$HOME/.local/share}/milpa"
target="${user_target}/repos"
if [[ "$MILPA_OPT_GLOBAL" ]]; then
  target="${MILPA_ROOT}/repos"
else
  # sometimes milpa gets installed as root, but each user
  # should have it's own user repo folder
  mkdir -p "$target"
fi

base="${MILPA_ARG_SOURCE%%/.milpa/*}"

function clean_filename() {
  local s="${1?need a string}"
  s="${s##*://}"
  s="${s##*@}"
  s="${s%%.git}"
  s="${s//[^[:alnum:]]/-}"
  s="${s//+(-)/-}"
  s="${s//--/-}"
  s="${s/#-}"
  s="${s/%-}"
  echo -n "${s}" | tr '[:upper:]' '[:lower:]'
}

function git_default_branch() {
  git ls-remote --symref "$1" HEAD | awk '/^ref/ {sub(".*/", "", $2); print $2}'
}

if [[ -d "$base/.milpa" ]]; then
  @milpa.log info "Local repository detected at $base, symlinking..."
  full_path="$(cd "$base" && pwd)" || @milpa.fail "Could not cd into local directory <$base>"
  repo_name="${full_path##*/}"
  repo_name="${repo_name#.}"
  new_repo="$target/$repo_name"
  [[ -d "$new_repo" ]] && @milpa.fail "A repo named $repo_name already exists at $new_repo"

  ln -sfv "$full_path/.milpa" "$new_repo" || @milpa.fail "Failed to symlink"
  @milpa.log success "Symlink created"
else
  @milpa.log info "git repository detected, cloning..."
  git_repos="${target%/repos}/clones"
  repo_name="$(clean_filename "$base")"
  repo_target="$git_repos/$repo_name"
  new_repo="$target/$repo_name"
  [[ -d "$repo_target" ]] && @milpa.fail "$repo_target already present"

  @milpa.log info "Cloning repository $base"
  mkdir -p "$git_repos" || @milpa.fail "could not create $git_repos"
  git init "$repo_target" || @milpa.fail "could not init $repo_target"
  cd "$repo_target" || @milpa.fail "Could not cd into $repo_target"
  git remote add origin "$base" || @milpa.fail "could not set origin to $base"
  git config core.sparsecheckout true || @milpa.fail "could not set config to do sparse checkouts"
  echo ".milpa/*" >> .git/info/sparse-checkout
  default_branch="$(git_default_branch "$base")" || @milpa.fail "Could not find default branch for $base"
  git pull --depth=1 origin "$default_branch" || @milpa.fail "Could not pull from $default_branch of $base"
  git switch "$default_branch"
  ln -sfv "$repo_target/.milpa" "$new_repo" || @milpa.fail "could not symlink to $target"
  @milpa.log success "Repo cloned"
fi

@milpa.log info "Running repo setup tasks"
if [[ -f "$new_repo/hooks/post-install.sh" ]]; then
  @milpa.log info "Running post-install hook"
  #shellcheck disable=1090,1091
  source "$new_repo/hooks/post-install.sh" || @milpa.fail "Could not run post-install hook to completion"
fi

@milpa.log complete "Repo installed at $new_repo"
