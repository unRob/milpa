#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

formatter="$MILPA_OPT_FORMAT"
if [[ "$formatter" == "auto" ]]; then
  if [[ "$CI" != "" ]]; then
    formatter="tap"
  else
    formatter="pretty"
  fi
fi

if [[ "${#MILPA_ARG_PATHS}" -eq 0 ]]; then
  MILPA_ARG_PATHS=( test/*.bats test/commands/**/*.bats )
fi

cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"
# shellcheck disable=2155
export TEST_MILPA_VERSION="$("$MILPA_ROOT/compa" __version 2>&1)"
@milpa.log info "Running integration tests"
export BATS_LIB_PATH="$MILPA_ROOT/test/_helpers"
env bats --timing --formatter "$formatter" "${MILPA_ARG_PATHS[@]}" || exit 2
@milpa.log complete "Integration tests passed"
