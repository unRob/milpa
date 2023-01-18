#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"
@milpa.log info "Running unit tests"
args=()
after_run=complete
if [[ "${MILPA_OPT_COVERAGE}" ]]; then
  args=( -coverprofile=test/coverage.out --coverpkg=./...)
fi
gotestsum --format short -- ./... "${args[@]}" || exit 2
@milpa.log "$after_run" "Unit tests passed"

[[ ! "${MILPA_OPT_COVERAGE}" ]] && exit
@milpa.log info "Building coverage report"
go tool cover -html=test/coverage.out -o test/coverage.html || @milpa.fail "could not build reports"
@milpa.log "$after_run" "Coverage report ready at test/coverage.html"
