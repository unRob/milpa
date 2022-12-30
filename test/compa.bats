#!/usr/bin/env bats
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
bats_load_library 'milpa'
_suite_setup

setup () {
  _common_setup
}

@test "compa prints version" {
  # compa only talks to stdout when talking to milpa
  # compa parses flags, so it should parse the version flag
  run -42 --keep-empty-lines --separate-stderr compa --version
  assert_equal "$output" ""
  assert_equal "$stderr" "$TEST_MILPA_VERSION"

  run -42 --keep-empty-lines --separate-stderr compa __version
  assert_equal "$output" ""
  assert_equal "$stderr" "$TEST_MILPA_VERSION"
}

@test "compa exits correctly on bad commands" {
  run -127 --separate-stderr compa bad-command
  assert_equal "$output" ""
  linecount=${#stderr_lines[@]}
  last_line=${stderr_lines[$(( linecount - 1))]}
  echo "$last_line"
  echo "${last_line}" | grep -m1 "Unknown subcommand bad-command"
}

@test "compa exits correctly on bad arguments" {
  run -64 --separate-stderr compa version --bad-flag
  assert_equal "$output" ""
  linecount=${#stderr_lines[@]}
  last_line=${stderr_lines[$(( linecount - 1))]}
  echo "$last_line"
  echo "${last_line}" | grep -m1 "unknown flag: --bad-flag"
}
