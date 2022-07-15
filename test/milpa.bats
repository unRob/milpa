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
load 'test/_helpers/setup.bash'
_suite_setup


setup () {
  load 'test/_helpers/bats-support/load.bash'
  load 'test/_helpers/bats-assert/load.bash'
  _common_setup
}

@test "milpa with no arguments shows help" {
  run milpa
  assert_failure 127
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
  run milpa
  assert_failure 78
}

@test "milpa includes global repos in MILPA_PATH" {
  run milpa debug-env MILPA_PATH
  assert_success
  assert_output "$(readlink -f "$MILPA_ROOT/.milpa"):$(readlink -f "$MILPA_ROOT/repos/test-suite")"
}

@test "milpa prepends user-supplied MILPA_PATH" {
  # path must have a milpa repo or it will be ignored!
  mkdir -pv "$BATS_SUITE_TMPDIR/somewhere/.milpa"
  export MILPA_PATH="$BATS_SUITE_TMPDIR/somewhere"
  run milpa debug-env MILPA_PATH
  assert_success
  assert_output "${BATS_SUITE_TMPDIR}/somewhere:$(readlink -f $MILPA_ROOT/.milpa):$(readlink -f $MILPA_ROOT/repos/test-suite)"
}
