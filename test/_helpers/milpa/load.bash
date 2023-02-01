#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

function _common_setup () {
  cd "$XDG_DATA_HOME" || exit 2
}

function _suite_setup() {
  bats_require_minimum_version 1.5.0
  bats_load_library 'bats-support'
  bats_load_library 'bats-assert'

  unset XDG_DATA_HOME MILPA_ROOT MILPA_PATH MILPA_PATH_PARSED DEBUG
  export XDG_DATA_HOME="${BATS_SUITE_TMPDIR//\/\///}"
  # shellcheck disable=2155
  export PROJECT_ROOT="$( cd "${BATS_TEST_FILENAME%%/test/*}" >/dev/null 2>&1 && pwd )"
  # make executables in src/ visible to PATH
  export MILPA_UPDATE_CHECK_DISABLED="yes"
  export MILPA_ROOT="$XDG_DATA_HOME/var/lib/milpa"
  export PATH="$PROJECT_ROOT:$PATH"
  export NO_COLOR=1
  export MILPA_HELP_STYLE="markdown"
  export milpa="$MILPA_ROOT/milpa"
  mkdir -p "$XDG_DATA_HOME"
  [[ -f "$XDG_DATA_HOME/setup-complete" ]] && return 0
  mkdir -pv "$MILPA_ROOT/repos"
  ln -sf "$PROJECT_ROOT/milpa" "$MILPA_ROOT/milpaa"
  ln -sf "$PROJECT_ROOT/.milpa" "$MILPA_ROOT/.milpa"
  mkdir -p "$XDG_DATA_HOME/milpa/repos"
  ln -sf "$PROJECT_ROOT/test/.milpa" "$MILPA_ROOT/repos/test-suite"
  touch "$XDG_DATA_HOME/setup-complete"
}

function _debug {
  echo "$@" >&3
}

function fixture() {
  echo "$PROJECT_ROOT/test/fixtures/$1"
}
