#!/usr/bin/env bats
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
bats_load_library 'milpa'
_suite_setup
bats_load_library 'bats-file'
export LOCAL_REPO="$XDG_DATA_HOME/.milpa"

setup() {
  _common_setup
  mkdir -p .milpa/commands
}

@test "itself create" {
  run milpa itself create
  assert_failure 64
}

@test "itself create something" {
  name="something-$(date -u '+%s')"
  unset MILPA_SILENT
  run milpa --verbose itself create "$name"
  assert_success
  assert_output --partial "$LOCAL_REPO/commands/$name.sh"
  assert_file_exist "$LOCAL_REPO/commands/$name.sh"
  assert_file_exist "$LOCAL_REPO/commands/$name.yaml"
  assert_file_not_executable "$LOCAL_REPO/commands/$name.sh"

  run milpa itself command-tree
  assert_output --partial "$name"

  run milpa "$name"
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

@test "itself create --kind executable something-executable" {
  run milpa itself create --kind executable something-executable
  assert_success
  assert_file_exist "$LOCAL_REPO/commands/something-executable"
  assert_file_exist "$LOCAL_REPO/commands/something-executable.yaml"
  assert_file_executable "$LOCAL_REPO/commands/something-executable"

  run milpa itself create --kind executable something-executable
  assert_failure 2
  assert_output --partial "Command already exists"

  echo '#!/usr/bin/env bash' >"$LOCAL_REPO/commands/something-executable"
  echo 'env | grep ^MILPA | sort' >>"$LOCAL_REPO/commands/something-executable"

  run milpa something-executable
  assert_output --partial "MILPA_COMMAND_KIND=executable"
  assert_output --partial "MILPA_COMMAND_NAME=milpa something-executable"
  refute_output --partial "MILPA_OPT_"
  refute_output --partial "MILPA_ARG_"
}

@test "itself create --kind zsh" {
  run milpa itself create --kind zsh zsh-thing
  assert_success
  assert_file_exist "$LOCAL_REPO/commands/zsh-thing.zsh"
  assert_file_exist "$LOCAL_REPO/commands/zsh-thing.yaml"

  run milpa itself create --kind zsh zsh-thing
  assert_failure 2
  assert_output --partial "Command already exists"

  echo 'env | grep ^MILPA | sort' >>"$LOCAL_REPO/commands/zsh-thing.zsh"

  run milpa zsh-thing
  assert_success
  assert_output --partial "MILPA_COMMAND_KIND=shell-script"
  assert_output --partial "MILPA_COMMAND_NAME=milpa zsh-thing"
  refute_output --partial "MILPA_OPT_"
  refute_output --partial "MILPA_ARG_"
}

@test "itself create --kind fish" {
  run milpa itself create --kind fish fish-thing
  assert_success
  assert_file_exist "$LOCAL_REPO/commands/fish-thing.fish"
  assert_file_exist "$LOCAL_REPO/commands/fish-thing.yaml"

  run milpa itself create --kind fish fish-thing
  assert_failure 2
  assert_output --partial "Command already exists"

  echo 'env | grep ^MILPA | sort' >>"$LOCAL_REPO/commands/fish-thing.fish"
  cat "$LOCAL_REPO/commands/fish-thing.fish"

  run milpa fish-thing
  assert_success
  assert_output --partial "MILPA_COMMAND_KIND=shell-script"
  assert_output --partial "MILPA_COMMAND_NAME=milpa fish-thing"
  refute_output --partial "MILPA_OPT_"
  refute_output --partial "MILPA_ARG_"
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
