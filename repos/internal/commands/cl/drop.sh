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

if [[ "$MILPA_OPT_KIND" != "" ]]; then
  git notes --ref="changelog/$MILPA_OPT_KIND" remove "$MILPA_ARG_REF" || @milpa.fail "Could not drop entries from $MILPA_ARG_REF"
  exit
fi

git for-each-ref --format="%(refname)" refs/notes/changelog |
  sed 's|.*changelog/||' |
  while read -r kind; do
    git notes --ref="changelog/$kind" remove "$MILPA_ARG_REF" 2>/dev/null && @milpa.log success "Dropped $kind entries"
  done
