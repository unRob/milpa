#!/usr/bin/env bash
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
