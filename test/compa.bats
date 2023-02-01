#!/usr/bin/env bats
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
bats_load_library 'milpa'
_suite_setup

setup () {
  _common_setup
}

@test "milpa prints version" {
  run --keep-empty-lines --separate-stderr milpa --version
  assert_equal "$output" "$TEST_MILPA_VERSION"
  assert_equal "$stderr" ""

  run --keep-empty-lines --separate-stderr milpa __version
  assert_equal "$stderr" "$TEST_MILPA_VERSION"
  assert_equal "$output" ""
}

@test "milpa exits correctly on bad MILPA_ROOT" {
  MILPA_ROOT="$BATS_TEST_FILENAME"
  run -78 milpa
}


@test "milpa exits correctly on bad commands" {
  run -127 --separate-stderr milpa bad-command
  assert_equal "$output" ""
  linecount=${#stderr_lines[@]}
  last_line=${stderr_lines[$(( linecount - 1))]}
  echo "$last_line"
  echo "${last_line}" | grep -m1 "Unknown subcommand bad-command"
}

@test "milpa exits correctly on bad arguments" {
  run -64 --separate-stderr milpa version --bad-flag
  assert_equal "$output" ""
  linecount=${#stderr_lines[@]}
  last_line=${stderr_lines[$(( linecount - 1))]}
  echo "$last_line"
  echo "${last_line}" | grep -m1 "unknown flag: --bad-flag"
}
