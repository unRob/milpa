#!usr/bin/env bats
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

setup_file () {
  load 'test/_helpers/setup.bash'
  _suite_setup
}

setup () {
  load 'test/_helpers/bats-support/load.bash'
  load 'test/_helpers/bats-assert/load.bash'
  cd "$XDG_DATA_HOME" || exit 2
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
