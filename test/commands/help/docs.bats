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
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/milpa/ > "readme.html"
  kill -SIGINT "$server"

  run diff "index.html" "$PROJECT_ROOT/test/fixtures/index.html"
  assert_success

  run diff "readme.html" "$PROJECT_ROOT/test/fixtures/readme.html"
  assert_success
}
