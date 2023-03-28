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
  # regenerate with
  # MILPA_HELP_STYLE=markdown MILPA_ROOT="$(pwd)" MILPA_PATH="$(pwd)/.milpa" MILPA_PATH_PARSED=true milpa help docs milpa | tee test/fixtures/readme.txt
  run diff -u -L "live" <(milpa help docs milpa) -L "fixture" <(cat "$(fixture readme.txt)")
  assert_success

  run diff -u -L "live" <(milpa help docs milpa environment) -L "fixture" <(cat "$(fixture environment.txt)")
  assert_success

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
  # regenerate with
  # MILPA_PATH="$(pwd)/.milpa:$(pwd)/test/.milpa" MILPA_PATH_PARSED=true milpa help docs --server
  # curl http://localhost:4242/ > test/fixtures/index.html
  # curl http://localhost:4242/help/docs/ > test/fixtures/docs.html
  # curl http://localhost:4242/help/docs/milpa/ > test/fixtures/readme.html
  # curl http://localhost:4242/help/docs/milpa/environment/ > test/fixtures/environment.html
  milpa help docs --server &
  server="$!"
  tries=0
  while ! curl --max-time 2 --fail --silent --show-error http://localhost:4242/ > "index.html"; do
    tries+=1
    if [[ "$tries" -gt 5 ]]; then
      exit 2
    fi

    echo "waiting for server to come up"
    sleep .1
  done
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/ > "docs.html"
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/milpa/ > "readme.html"
  curl --max-time 2 --fail --silent --show-error http://localhost:4242/help/docs/milpa/environment/ > "environment.html"
  run -22 curl --max-time 2 --fail --silent --show-error http://localhost:4242/does-not-exist

  kill -SIGINT "$server"

  run diff -u -L fixture "$(fixture index.html)" "index.html"
  assert_success

  run diff -u -L fixture "$(fixture readme.html)" "readme.html"
  assert_success

  run diff -u -L fixture "$(fixture environment.html)" "environment.html"
  assert_success

  run diff -u -L fixture "$(fixture docs.html)" "docs.html"
  assert_success
}
