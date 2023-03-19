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

@test "itself command-tree (autocomplete)" {
  run compa __command_tree --format autocomplete
  assert_success
  assert_output "debug-env
itself"
}

@test "itself command-tree --output json" {
  milpa itself command-tree --depth 1 --output json > cmdtree.json

  run jq -e '(.children | length) == 2' cmdtree.json
  assert_success

  run jq -r '.children | first | .command.path | join(" ")' cmdtree.json
  assert_success
  assert_output "milpa debug-env"

  run jq -r '.children | last | .command.path | join(" ")' cmdtree.json
  assert_success
  assert_output "milpa itself"

  run jq -r '.children | last | .command.meta.kind' cmdtree.json
  assert_success
  assert_output "virtual"
}

@test "itself command-tree --output yaml" {
  milpa itself command-tree --depth 1 --output yaml | yq  -o json > cmdtree.yaml

  run jq -e '(.children | length) == 2' cmdtree.yaml
  assert_success

  run jq -r '.children | first | .command.path | join(" ")' cmdtree.json
  assert_success
  assert_output "milpa debug-env"

  run jq -r '.children | last | .command.path | join(" ")' cmdtree.json
  assert_success
  assert_output "milpa itself"

  run jq -r '.children | last | .command.meta.kind' cmdtree.json
  assert_success
  assert_output "virtual"
}
