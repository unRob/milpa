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
last_release="$(git describe --abbrev=0 --exclude='*-*' --tags 2>/dev/null || git rev-list --max-parents=0 HEAD)"

main_release_pattern="^[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+$"
if [[ $this_release =~ $main_release_pattern ]]; then
  # main release
  block="$this_release"
  header="## [$this_release](https://github.com/unRob/milpa/releases/tag/$this_release) - $(date -u "+%Y-%m-%d")"
  # remove upcoming from changelog
  drop_block="/<!-- upcoming -->/,/<!-- upcoming !-->/d;"
else
  block="upcoming"
  drop_block=""
  header="## [Upcoming](https://github.com/unRob/milpa/compare/$last_release...HEAD)"
fi

@milpa.log info "Generating release notes for $last_release...$this_release"
cd "$MILPA_ROOT" || @milpa.fail "unknown milpa_root"
changelog-build \
  -last-release="${last_release}" \
  -this-release="${MILPA_CL_VERSION:-$this_release}" \
  -entries-dir="internal/.changelog/entries" \
  -changelog-template="internal/.changelog/changelog.tmpl" \
  -note-template="internal/.changelog/note.tmpl" > "dist/changes.md" || @milpa.fail "Could not generate release notes"

if [[ ! "$MILPA_OPT_SKIP_UPDATE" ]]; then
  previousNotes=$(sed "1,2d; $drop_block /<!-- $block -->/,/<!-- $block !-->/d; /./,\$!d" CHANGELOG.md)

  if [[ "$previousNotes" != "" ]]; then
    previousNotes=$'\n'$'\n'"$previousNotes"
  fi

  cat > CHANGELOG.md <<MD
# Changelog

<!-- $block -->
$header
$(cat "$MILPA_ROOT/dist/changes.md")
<!-- $block !-->${previousNotes}
MD
fi

if [[ "$MILPA_OPT_OUTPUT" ]]; then
  cat "dist/changes.md"
fi

rm dist/changes.md
