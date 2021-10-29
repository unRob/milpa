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

setup_file() {
  load 'test/_helpers/setup.bash'
  _suite_setup
}

setup() {
  load 'test/_helpers/bats-support/load.bash'
  load 'test/_helpers/bats-assert/load.bash'
  load 'test/_helpers/bats-file/load.bash'
  cd "$BATS_TEST_TMPDIR"
  mkdir -p .milpa/commands
}


@test "itself create" {
  run milpa itself create
  assert_failure 64
}

@test "itself create something" {
  run milpa itself create something
  assert_success
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something.sh"
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something.yaml"
  assert_file_not_executable "$BATS_TEST_TMPDIR/.milpa/commands/something.sh"
}

@test "itself create something deeply nested" {
  run milpa itself create something deeply nested
  assert_success
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something/deeply/nested.sh"
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something/deeply/nested.yaml"
}

@test "itself create something --executable" {
  run milpa itself create something --executable
  assert_success
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something"
  assert_file_exist "$BATS_TEST_TMPDIR/.milpa/commands/something.yaml"
  assert_file_executable "$BATS_TEST_TMPDIR/.milpa/commands/something"
}

@test "itself create something --repo somewhere-else" {
  run milpa itself create something --repo "$BATS_TEST_TMPDIR/somewhere-else"
  assert_success
  assert_file_exist "$BATS_TEST_TMPDIR/somewhere-else/.milpa/commands/something.sh"
  assert_file_exist "$BATS_TEST_TMPDIR/somewhere-else/.milpa/commands/something.yaml"
}
