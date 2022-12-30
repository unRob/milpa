#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
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
