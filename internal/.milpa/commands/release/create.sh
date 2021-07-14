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

if [[ "$MILPA_ARG_INCREMENT" == "" ]]; then
  notes="$(milpa release notes --output --skip-update --silent)"
  case "$notes" in
    *"### ✂️ Breaking Changes"*) MILPA_ARG_INCREMENT="major" ;;
    *"### 🌱 Features"*|*"### 🌺 Improvements"*) MILPA_ARG_INCREMENT="minor" ;;
    *) MILPA_ARG_INCREMENT="patch";
  esac

  @milpa.log info "Auto-detected semver increment of type: $(@milpa.fmt bold $MILPA_ARG_INCREMENT)"
fi


current_branch=$(git rev-parse --abbrev-ref HEAD)
[[ "$current_branch" != "main" ]] && @milpa.fail "Refusing to release on branch <$current_branch>"
[[ -n "$(git status --porcelain)" ]] && @milpa.fail "Git tree is messy, won't continue"

# get the latest tag, ignoring any pre-releases
# by default current version is 0.-1.-1, since any initial release will include features
# and thus be a minor release
current="$(git describe --abbrev=0 --exclude='*-*' --tags 2>/dev/null || echo "0.-1.-1")"

next=$(next_semver "$MILPA_ARG_INCREMENT" "$current")

if [[ "$MILPA_OPT_PRE" ]]; then
  # pre releases might update previous ones, look for them
  pre_current=$(git describe --abrev=0 --match="$next-$MILPA_OPT_PRE.*" --tags 2>/dev/null || echo "$current-$MILPA_OPT_PRE.-1")
  build=${pre_current##*.}
  next="$next-$MILPA_OPT_PRE.$(( build + 1 ))"
  notes="$(milpa release notes --output --skip-update --silent)"
fi

@milpa.log info "Creating release with version $(@milpa.fmt inverted "$next")"

if [[ ! "$MILPA_OPT_PRE" ]]; then
  # mainline releases need updated changelogs
  @milpa.log info "Updating CHANGELOG.md"
  notes=$(MILPA_CL_VERSION="HEAD" milpa release notes "$next" --output)
  {
    git add "CHANGELOG.md" && git commit -m "Update changelog for release $next" && git push;
  } || @milpa.fail "Could not commit CHANGELOG update"
fi

@milpa.log info "Creating tag and pushing"
git tag "$next" -m "$notes" || @milpa.fail "Could not create tag $next"
git push origin "$next" || @milpa.fail "Could not push tag $next"
