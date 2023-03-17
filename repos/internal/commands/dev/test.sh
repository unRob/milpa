#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

set -e
if [[ "$MILPA_OPT_COVERAGE" ]]; then
  milpa dev build --coverage
  milpa dev test unit --coverage
  milpa dev test integration --coverage
  milpa dev test coverage-report
else
  milpa dev test unit
  milpa dev test integration
fi
