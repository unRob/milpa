#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"
@milpa.log info "Running unit tests"
args=()
after_run=complete
if [[ "${MILPA_OPT_COVERAGE}" ]]; then
  cover_dir="$MILPA_ROOT/test/coverage/unit"
  rm -rf "$cover_dir"
  mkdir -p "$cover_dir"
  args=( -test.gocoverdir="$cover_dir" --coverpkg=./...)
fi
gotestsum --format short -- ./... "${args[@]}" || exit 2
@milpa.log "$after_run" "Unit tests passed"
