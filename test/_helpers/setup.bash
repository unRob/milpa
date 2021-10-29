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

function _suite_setup() {
  export XDG_DATA_HOME="${BATS_SUITE_TMPDIR//\/\///}/home"
  # shellcheck disable=2155
  export PROJECT_ROOT="$( cd "${BATS_TEST_FILENAME%%/test/*}" >/dev/null 2>&1 && pwd )"
  # make executables in src/ visible to PATH
  export MILPA_DISABLE_UPDATE_CHECKS="yes"
  export MILPA_ROOT="$XDG_DATA_HOME/var/lib/milpa"
  export PATH="$PROJECT_ROOT:$PATH"
  export NO_COLOR=1
  export MILPA_PLAIN_HELP=enabled

  [[ -f "$XDG_DATA_HOME/setup-complete" ]] && return 0
  mkdir -pv "$MILPA_ROOT/repos"
  ln -sfv "$PROJECT_ROOT/milpa" "$MILPA_ROOT/milpa"
  ln -sfv "$PROJECT_ROOT/compa" "$MILPA_ROOT/compa"
  ln -sfv "$PROJECT_ROOT/.milpa" "$MILPA_ROOT/.milpa"
  mkdir -pv "$XDG_DATA_HOME/.local/share/milpa/repos"
  ln -sfv "$PROJECT_ROOT/test/.milpa" "$MILPA_ROOT/repos/test-suite"
  touch "$XDG_DATA_HOME/setup-complete"
}

function _debug {
  echo "$@" >&3
}
