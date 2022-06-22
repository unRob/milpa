#!/usr/bin/env bats
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
export LOCAL_REPO="$XDG_DATA_HOME/.milpa"

setup() {
  load 'test/_helpers/bats-support/load.bash'
  load 'test/_helpers/bats-assert/load.bash'
  load 'test/_helpers/bats-file/load.bash'
  _common_setup
  mkdir -p .milpa/commands
}

@test "itself create" {
  run milpa itself create
  assert_failure 64
}

@test "itself create something" {
  run milpa --verbose itself create something
  assert_success
  assert_output --partial "$LOCAL_REPO/commands/something.sh"
  assert_file_exist "$LOCAL_REPO/commands/something.sh"
  assert_file_exist "$LOCAL_REPO/commands/something.yaml"
  assert_file_not_executable "$LOCAL_REPO/commands/something.sh"

  run milpa itself command-tree
  assert_output --partial "something"

  run milpa something
  assert_success
  assert_output ""
}

@test "itself create something-else deeply nested" {
  run milpa itself create something-else deeply nested
  assert_success
  assert_file_exist "$LOCAL_REPO/commands/something-else/deeply/nested.sh"
  assert_file_exist "$LOCAL_REPO/commands/something-else/deeply/nested.yaml"

  run milpa something-else deeply nested
  assert_success
  assert_output ""
}

@test "itself create something-executable --executable" {
  run milpa itself create something-executable --executable
  assert_success
  assert_file_exist "$LOCAL_REPO/commands/something-executable"
  assert_file_exist "$LOCAL_REPO/commands/something-executable.yaml"
  assert_file_executable "$LOCAL_REPO/commands/something-executable"

  run milpa itself create something-executable --executable
  assert_failure 2
  assert_output --partial "Command already exists"
}

@test "itself create something-elsewhere --repo somewhere-else" {
  run milpa itself create something-elsewhere --repo "$BATS_TEST_TMPDIR/somewhere-else/.milpa"
  assert_success
  assert_file_exist "$BATS_TEST_TMPDIR/somewhere-else/.milpa/commands/something-elsewhere.sh"
  assert_file_exist "$BATS_TEST_TMPDIR/somewhere-else/.milpa/commands/something-elsewhere.yaml"

  mkdir -p "$BATS_TEST_TMPDIR/somewhere-else/.milpa"
  cd "$BATS_TEST_TMPDIR/somewhere-else"
  run milpa something-elsewhere
  assert_success
  assert_output ""
}
