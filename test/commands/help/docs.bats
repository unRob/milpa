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

@test "itself docs" {
  run milpa help docs milpa
  assert_success
  assert_output "$(cat "$(fixture readme.txt)")"

  run milpa help docs milpa environment
  assert_success
  assert_output "$(cat "$(fixture environment.txt)")"

  run milpa __complete help docs ""
  assert_success
  assert_output "milpa
_activeHelp_ The topic to show docs for
:4
Completion ended with directive: ShellCompDirectiveNoFileComp"

  run milpa __complete help docs milpa "i"
  assert_success
  assert_output "index
internals
:4
Completion ended with directive: ShellCompDirectiveNoFileComp"
}


@test "itself docs --server" {
  milpa help docs --server &
  server="$!"
  tries=0
  while ! curl --max-time 2 --fail --silent --show-error http://localhost:4242/ > "index.html"; do
    tries+=1
    if [[ "$tries" -gt 5 ]]; then
      exit 2
    fi

    echo "waiting for server to come up"
    sleep 1
  done
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/ > "docs.html"
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/milpa/ > "readme.html"
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/milpa/environment/ > "environment.html"
  run -22 curl --max-time 2 --fail --silent --show-error http://localhost:4242/does-not-exist

  kill -SIGINT "$server"

  run diff "index.html" "$(fixture index.html)"
  assert_success

  run diff "readme.html" "$(fixture readme.html)"
  assert_success

  run diff "environment.html" "$(fixture environment.html)"
  assert_success

  run diff "docs.html" "$(fixture docs.html)"
  assert_success
}
