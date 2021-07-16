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

this_release="$MILPA_ARG_VERSION"
last_release="$(git describe --abbrev=0 --exclude='*-*' --exclude="$this_release" --tags "$this_release" 2>/dev/null)"
if [[ "$last_release" == "" ]]; then
  initial_commit=$(git rev-list --max-parents=0 HEAD)
fi

[[ "${MILPA_OPT_HEADER_OFFSET}" ]] && hl="$(printf -- "#%.0s" $(seq 1 "${MILPA_OPT_HEADER_OFFSET}"))"

function entryHeader() {
  case "$1" in
    breaking-change) echo "✂️ Breaking changes" ;;
    bug) echo "💦 Bug fixes" ;;
    feature) echo "🌱 Features" ;;
    improvement) echo "🌺 Improvements" ;;
    deprecation) echo "🍂 Deprecations" ;;
    note) echo "🧑🏽‍🌾 Notes" ;;
    *) @milpa.log warning "Unknown changelog kind: <$1>"; echo "🧑🏽‍🌾 Notes" ;;
  esac
}

# get all types of changelog entries
kinds=$(git for-each-ref --format="%(refname)" refs/notes/changelog | sed 's|.*changelog/||')
# get all commits of this range
if [[ "$last_release" == "" ]]; then
  commits=$(git rev-list "$this_release")
else
  commits=$(git rev-list "$(git rev-list -n 1 "$last_release")..$this_release")
fi
notes=""

if [[ "$commits" == "" ]]; then
  @milpa.fail "No new commits since ${last_release:-$initial_commit}"
fi

for kind in $kinds; do
  notesCommits=$(git notes --ref="changelog/$kind" list | cut -d' ' -f2)
  [[ "$notesCommits" == "" ]] && continue
  commitsWithNotes=$(grep -f <(printf '%s' "$notesCommits") <(printf "%s" "$commits") 2>/dev/null)

  [[ "$commitsWithNotes" == "" ]] && continue
  for commit in $commitsWithNotes; do
    prefix="- (${commit:0:6}) "
    notes+="$hl## $(entryHeader "$kind")"$'\n'$'\n'
    notes+="$(git notes --ref="changelog/$kind" show "$commit" | sed "s/🌽#🌽/$prefix/g")"
    notes+=$'\n'$'\n'
  done
done

if [[ "$notes" == "" ]]; then
  @milpa.fail "No release notes for $this_release (since ${last_release:-initial commit})"
fi

# we actually mean to strip newlines from the end
# shellcheck disable=2059
printf "${notes}"
