#!/usr/bin/env bats
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
bats_load_library 'milpa'
_suite_setup

setup () {
  _common_setup
}

@test "milpa with no arguments shows help" {
  run -127 milpa
  assert_output --regexp "## Usage"
}

@test "milpa help exits cleanly" {
  run milpa help
  assert_success
  assert_output --regexp "## Usage"

  run milpa --help
  assert_success
  assert_output --regexp "## Usage"
}

@test "milpa with bad MILPA_ROOT" {
  MILPA_ROOT="$BATS_TEST_FILENAME"
  run -78 milpa
}

@test "milpa includes global repos in MILPA_PATH" {
  run milpa debug-env MILPA_PATH
  assert_success
  assert_output "$(readlink -f "$MILPA_ROOT/.milpa"):$(readlink -f "$MILPA_ROOT/repos/test-suite")"
}

@test "milpa prepends user-supplied MILPA_PATH" {
  # path must have a milpa repo or it will be ignored!
  mkdir -pv "${BATS_SUITE_TMPDIR}/somewhere/.milpa"
  export MILPA_PATH="${BATS_SUITE_TMPDIR}/somewhere"
  run milpa debug-env MILPA_PATH
  assert_success
  assert_output "${BATS_SUITE_TMPDIR//\/\///}/somewhere/.milpa:$(readlink -f $MILPA_ROOT/.milpa):$(readlink -f $MILPA_ROOT/repos/test-suite)"
}
