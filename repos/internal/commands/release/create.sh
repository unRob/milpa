#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
@milpa.load_util user-input

current_branch=$(git rev-parse --abbrev-ref HEAD)
[[ "$current_branch" != "main" ]] && @milpa.fail "Refusing to release on branch <$current_branch>"
[[ -n "$(git status --porcelain)" ]] && @milpa.fail "Git tree is messy, won't continue"

function next_semver() {
  local components
  IFS="." read -r -a components <<< "${2}"
  following=""
  case "$1" in
    major ) following="$((components[0]+1)).0.0" ;;
    minor ) following="${components[0]}.$((components[1]+1)).0" ;;
    patch ) following="${components[0]}.${components[1]}.$((components[2]+1))" ;;
    *) @milpa.fail "unknown increment type: <$1>"
  esac

  echo "$following"
}

increment="$(milpa cl next-version)"
@milpa.log info "Auto-detected semver increment of type: $(@milpa.fmt bold "$increment")"

# get the latest tag, ignoring any pre-releases
# by default current version is 0.-1.-1, since any initial release will include features
# and thus be a minor release
current="$(git describe --abbrev=0 --exclude='*-*' --tags 2>/dev/null || echo "0.-1.-1")"

next=$(next_semver "$increment" "$current")

if [[ "$MILPA_OPT_PRE" ]]; then
  # pre releases might update previous ones, look for them
  pre_current=$(git describe --abbrev=0 --match="$next-$MILPA_OPT_PRE.*" --tags 2>/dev/null || echo "$current-$MILPA_OPT_PRE.-1")
  build=${pre_current##*.}
  next="$next-$MILPA_OPT_PRE.$(( build + 1 ))"
fi

@milpa.log info "Creating release with version $(@milpa.fmt inverted "$next")"
@milpa.confirm "Proceed with release?" || @milpa.fail "Refusing to continue, got <$REPLY>"
@milpa.log success "Continuing with release"


if [[ ! "$MILPA_OPT_PRE" ]]; then
  # mainline releases need updated changelogs
  @milpa.log info "Updating CHANGELOG.md:"
  set -o pipefail
  milpa cl update "$next" | sed '/^---/q' || @milpa.fail "Could not generate changelog"
  [[ -n "$(git status --porcelain)" ]] || @milpa.fail "Changelog was not updated"
  {
    git add "CHANGELOG.md" &&
    git commit --no-verify -m "Update changelog for release $next" &&
    git push;
  } || @milpa.fail "Could not commit and push CHANGELOG update"
fi

@milpa.log info "Creating tag and pushing"
git config core.commentChar ';'
git tag "$next" -m "$(milpa cl show)" || @milpa.fail "Could not create tag $next"
git config core.commentChar auto
git push origin "$next" || @milpa.fail "Could not push tag $next"

@milpa.log complete "Release created and pushed to origin!"
