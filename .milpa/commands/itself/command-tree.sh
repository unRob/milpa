#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

set -o pipefail

if [[ "$MILPA_ARG_PREFIX" ]]; then
  @milpa.log info "looking for commands with prefix <$MILPA_ARG_PREFIX>"
else
  @milpa.log info "looking for all known commands"
fi

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
