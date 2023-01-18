#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

@milpa.load_util repo

if [[ "$MILPA_OPT_REPO" == "" ]]; then
  repo_path="$(@milpa.repo.current_path)"
else
  repo_path="$MILPA_OPT_REPO"
fi

milpa="$repo_path/.milpa"
joinedName="${MILPA_ARG_NAME[*]}"
path="$milpa/docs/${joinedName// /\/}.md"
@milpa.log info "Creating doc for $(@milpa.fmt bold "${MILPA_ARG_NAME[*]}") at $path"
mkdir -p "$(dirname "$path")"

cat > "$path" <<MD
---
title: $joinedName
---

This document talks about $joinedName
MD

@milpa.log complete "doc $(@milpa.fmt bold "${MILPA_ARG_NAME[*]}") created"
[[ "$MILPA_OPT_OPEN" ]] && $EDITOR "$path"
