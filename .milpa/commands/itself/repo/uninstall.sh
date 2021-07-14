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

@milpa.log info "Removing $(@milpa.fmt bold "$MILPA_ARG_PATH")"
if [[ -f "$MILPA_ARG_PATH/hooks/post-uninstall.sh" ]]; then
  @milpa.log info "Running post-uninstall hook"
  # run in a subshell so we don't care if it uninstall hook does weird stuff
  (
    #shellcheck disable=1090,1091
    source "$MILPA_ARG_PATH/hooks/post-uninstall.sh"
  ) || @milpa.log warning "Could not run post-uninstall hook to completion"
fi

rm -rf "$MILPA_ARG_PATH"
