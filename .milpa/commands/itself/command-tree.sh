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
# shellcheck disable=2048,2086
"$MILPA_COMPA" __inspect ${MILPA_ARG_PREFIX[*]} |
  jq -r 'def toIndentedTree($indent):
    reduce .[] as $leaf ([];
      . + [ ([$indent, ($leaf.command.Meta.Name | last), $leaf.command.Summary] | join("¬")) ] + (if $leaf.children then ($leaf.children | toIndentedTree($indent + "  ")) else [] end)
    );
  .children | toIndentedTree("")[]' |
  while IFS="¬" read -r offset name description; do
    prefix="$offset"
    suffix=" - $description"
    if [[ "$MILPA_OPT_NAME_ONLY" ]]; then
      prefix=""
      suffix=""
    fi

    if [[ "$(( ${#offset} / 2 ))" -lt "$MILPA_OPT_DEPTH" ]]; then
      echo "${prefix}$(@milpa.fmt bold "$name")$suffix"
    fi
  done

