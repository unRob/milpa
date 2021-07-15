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

{
cat <<MD
# Changelog

<!-- this changelog is generated with \`milpa cl update\` -->
Milpa follows the [semver 2.0.0](https://semver.org/spec/v2.0.0.html) specification.

MD

# print upcoming, and ignore if there's nothing to do
if notes=$(milpa cl show --header-offset 1 2>/dev/null); then
  @milpa.log info "Adding upcoming entries"
  previous="$(git describe --abbrev=0 --exclude='*-*' --exclude="$MILPA_ARG_VERSION" --tags 2>/dev/null)"
  if [[ "$previous" == "" ]]; then
    previous=$(git rev-list --max-parents=0 HEAD)
  fi

  if [[ "$MILPA_ARG_VERSION" == "HEAD" ]]; then
    echo "# [Upcoming](https://github.com/unRob/milpa/compare/${previous}...HEAD)"
  else
    echo "# [$MILPA_ARG_VERSION](https://github.com/unRob/milpa/releases/tag/$MILPA_ARG_VERSION) - $(date -u "+%Y-%m-%d")"
  fi

  echo
  echo "$notes"
  echo
  echo "---"
fi

git tag --sort=-taggerdate | grep '^\d\+\.\d\+\.\d\+$' |
  while read -r tag; do
    printf '\n## [%s](https://github.com/unRob/milpa/releases/tag/%s) - ' "$tag" "$tag"
    git tag --points-at "$tag" \
      --format="%(taggerdate:short)"$'\n'$'\n'"%(contents)" |
      sed 's/^## /### /g'
    echo "---"
  done
} | tee "$MILPA_ROOT/CHANGELOG.md" | less -FIRX
