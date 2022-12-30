#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function @milpa.repo.current_path () {
  repo_path="$(pwd)"
  while [[ ! -d "$repo_path/.milpa" ]]; do
    [[ "$repo_path" == "/" ]] && return 2
    repo_path=$(dirname "$repo_path")
  done
  echo "$repo_path"
}
