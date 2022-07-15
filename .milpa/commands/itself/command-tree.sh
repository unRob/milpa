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

set -o pipefail

@milpa.log info "looking for commands with prefix <$MILPA_ARG_PREFIX>"

function get_tree () {
  "$MILPA_COMPA" __inspect \
    --depth "$MILPA_OPT_DEPTH" \
    --format "$1" \
    "${2+--template=}${2:-}" \
    "${MILPA_ARG_PREFIX[@]}"
}

if [[ "$MILPA_OPT_OUTPUT" =~ ^(yaml|json)$ ]]; then
  get_tree "$MILPA_OPT_OUTPUT" || @milpa.fail "Could not load tree"
  exit
fi

if [[ "$MILPA_OPT_TEMPLATE" != "" ]]; then
  get_tree text "$MILPA_OPT_TEMPLATE" || @milpa.fail "Could not load tree"
  exit
fi

initialDepth="${#MILPA_ARG_PREFIX[@]}"
while IFS='¬' read -r depth name description; do
  depth=$(( depth - 1 - initialDepth ))
  indent=""
  if [[ "$depth" -gt 0 ]]; then
    indent="$(printf -- ' %.0s' $(seq 0 $depth))"
  fi
  echo "$indent$(@milpa.fmt bold "$name") - $description"
done < <(get_tree text "{{ len .Meta.Name }}¬{{ .Name }}¬{{ .Summary }}"$'\n')
