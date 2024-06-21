#!/usr/bin/env bats
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
bats_load_library 'milpa'
_suite_setup

setup () {
  _common_setup
}

@test "milpa with bad MILPA_ROOT" {
  MILPA_ROOT="$BATS_TEST_FILENAME"
  run -78 milpa
}

@test "milpa errors on bad configs" {
  repo="${BATS_SUITE_TMPDIR}/bad-config/.milpa"
  mkdir -pv "$repo/commands"
  echo "summary:"$'\n'"  int: - 1 :a\\" > "$repo/commands/bad-command.yaml"
  cat "$repo/commands/bad-command.yaml"
  echo "#!/usr/bin/env bash\necho not bad" > "$repo/commands/bad-command.sh"
  export MILPA_PATH="${BATS_SUITE_TMPDIR}/bad-config"
  run milpa bad-command
  assert_output --regexp "Run \`milpa itself doctor\` to diagnose your installed commands."
  assert_output --regexp "milpa bad-command"
  assert_output --regexp "$repo/commands/bad-command.yaml"
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
