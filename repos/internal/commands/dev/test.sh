#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

set -e
if [[ "$MILPA_OPT_COVERAGE" ]]; then
  rm -rf test/coverage
  milpa dev build --coverage
  milpa dev test unit --coverage
  milpa dev test integration --coverage
  mkdir -p "$MILPA_ROOT/test/coverage/doctor"
  MILPA_PATH="$(pwd)/.milpa" MILPA_PATH_PARSED=true GOCOVERDIR="$MILPA_ROOT/test/coverage/doctor" milpa itself doctor
  milpa dev test coverage-report
  milpa dev build
else
  milpa dev test unit
  milpa dev test integration
  milpa itself doctor
fi
