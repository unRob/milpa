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
