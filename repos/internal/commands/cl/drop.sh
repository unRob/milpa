#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

if [[ "$MILPA_OPT_KIND" != "" ]]; then
  git notes --ref="changelog/$MILPA_OPT_KIND" remove "$MILPA_ARG_REF" || @milpa.fail "Could not drop entries from $MILPA_ARG_REF"
  exit
fi

git for-each-ref --format="%(refname)" refs/notes/changelog |
  sed 's|.*changelog/||' |
  while read -r kind; do
    git notes --ref="changelog/$kind" remove "$MILPA_ARG_REF" 2>/dev/null && @milpa.log success "Dropped $kind entries"
  done
