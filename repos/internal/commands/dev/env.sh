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

@milpa.load_util shell
@milpa.log info "Exporting shell variables"

root="${MILPA_COMMAND_REPO%%/repos/internal*}"
@milpa.shell.export "MILPA_ROOT" "$root"

# shellcheck disable=2049
if [[ ! "$PATH" = "$root"* ]]; then
  @milpa.shell.prepend_path "$root"
fi
