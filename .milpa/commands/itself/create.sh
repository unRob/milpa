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

@milpa.load_util repo

if [[ "$MILPA_OPT_REPO" == "" ]]; then
  if ! repo_path="$(@milpa.repo.current_path)/.milpa"; then
    @milpa.log warning "No milpa repo detected, creating one at $(pwd)"
    repo_path="$(pwd)/.milpa"
  fi
else
  repo_path="$MILPA_OPT_REPO"
fi

command_name=${MILPA_ARG_NAME[*]}
command_path="$repo_path/commands/${command_name// //}"

[[ "${MILPA_OPT_EXECUTABLE}" ]] || command_path="$command_path.sh"

[[ -f "$command_path" ]] && @milpa.fail "Command already exists at $command_path"

@milpa.log info "Creating command $(@milpa.fmt bold "${command_name}") at $command_path"
mkdir -p "$(dirname "$command_path")"

if [[ "${MILPA_OPT_EXECUTABLE}" ]]; then
  touch "$command_path" || @milpa.fail "could not create $command_path"
  chmod +x "$command_path"
else
  echo "#!/usr/bin/env bash" >> "$command_path" || @milpa.fail "could not create $command_path"
fi

# shellcheck disable=2001
cat > "${command_path%.sh}.yaml" <<YAML
# see \`milpa help docs command spec\` for all the options
summary: ${MILPA_OPT_SUMMARY}
description: |
$(sed 's/^/  /' <<<"${MILPA_OPT_DESCRIPTION}")
# arguments:
#   - name: something
#     description: Sets SOMETHING
#     required: true
#     variadic: false
#     values:
#       script: whoami
# options:
#   option:
#     description: sets OPTION
#     default: fourty-two
#     values:
#       static: [one, two, fourty-two]
YAML

@milpa.log complete "$(@milpa.fmt bold "${command_name}") created"
[[ "$MILPA_OPT_OPEN" ]] && $EDITOR "$command_path"

exit 0
