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

this_release="HEAD"
last_release="$(git describe --abbrev=0 --exclude='*-*' --tags 2>/dev/null)"

if [[ "$last_release" == "" ]]; then
  commits=$(git rev-list "$this_release") || @milpa.fail "Could not get list of commits"
else
  commits=$(git rev-list "$(git rev-list -n 1 "$last_release")..$this_release") || @milpa.fail "Could not get list of commits"
fi

function has_entries_of_kind() {
  notesCommits=$(git notes --ref="changelog/$1" list | cut -d' ' -f2) || return 1
  grep -f <(printf '%s' "$notesCommits") <(printf "%s" "$commits") >/dev/null 2>&1
}

if has_entries_of_kind breaking-change; then
  echo "major"
elif has_entries_of_kind feature || has_entries_of_kind improvement; then
  echo "minor"
else
  echo "patch"
fi
