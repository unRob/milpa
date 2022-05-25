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

@milpa.log info "Linting shell scripts"
find .milpa repos/internal/commands -name '*.sh' -exec shellcheck {} \+  || @milpa.fail "could not lint commands"
shellcheck milpa bootstrap.sh test/_helpers/*.bash || @milpa.fail "could not lint helper files"
@milpa.log complete "Shell files are up to spec"
