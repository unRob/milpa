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
@milpa.load_util user-input

entry="${MILPA_ARG_MESSAGE[*]}"
if [[ "$entry" == "" || "$entry" == "-" ]]; then
  if [[ -t 1 ]]; then
    entry=$(@milpa.ask "Enter the message for this <$MILPA_ARG_KIND>:")
  else
    entry=$(cat)
  fi
fi

if [[ "$entry" == "" ]]; then
  @milpa.fail "No entry message supplied"
fi

printf '🌽#🌽%s' "$entry" |
  git notes --ref "changelog/$MILPA_ARG_KIND" \
    append --file - "$MILPA_OPT_REF" ||
      @milpa.fail "Could not add git note to $MILPA_OPT_REF"
