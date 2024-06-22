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

@test "milpa with no arguments shows help" {
  run -127 milpa
  assert_output --regexp "## Usage"
}

@test "milpa help exits cleanly" {
  # shows help on the help command
  run milpa help
  assert_success
  assert_output --regexp "milpa help SUBCOMMAND"
  assert_output --regexp "Display usage information for any command"

  # shows help on milpa itself
  run milpa --help
  assert_success
  assert_output --regexp "Runs commands found in .milpa folders"

  # shows help on an existing sub command
  run milpa help itself create
  assert_success
  assert_output --regexp '`milpa itself create'

  # bad sub-command shows help of parent
  run -127 milpa help itself typotypo
  assert_failure 127
  assert_output --regexp '`milpa itself SUBCOMMAND'
  assert_output --regexp 'Unknown help topic \"typotypo\" for milpa itself'
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

@test "milpa completes recursively" {
  # path must have a milpa repo or it will be ignored!
  run milpa __complete debug-env --completion-test ""
  assert_success
  assert_output "$MILPA_ROOT
_activeHelp_ tests if milpa can call itself during completion
:0
Completion ended with directive: ShellCompDirectiveDefault"
}

@test "milpa renders properly with/out truecolor enabled" {
  sameError='\\E\[1;41;37m ERROR \\E\[22;0;0m Unknown subcommand bad-command'
  run bash -c 'printf '%q' "$(COLORTERM="terminal.app" NO_COLOR="" COLOR=always MILPA_HELP_STYLE=auto milpa bad-command 2>&1)"'
  assert_success
  assert_output --regexp "$sameError"
  assert_output --regexp '\\E\[38;5;193;1m\\E\[0m\\E\[38;5;193;1m\\E\[0m\\E\[38;5;193;1m## \\E\[0m\\E\[38;5;193;1mUsage\\E\[0m'

  run bash -c 'printf '%q' "$(COLORTERM="truecolor" NO_COLOR="" COLOR=always MILPA_HELP_STYLE=auto milpa bad-command 2>&1)"'
  assert_success
  assert_output --regexp "$sameError"
  assert_output --regexp '\\E\[38;2;192;227;147;1m\\E\[0m\\E\[38;2;192;227;147;1m\\E\[0m\\E\[38;2;192;227;147;1m## \\E\[0m\\E\[38;2;192;227;147;1mUsage\\E\[0m'
}
