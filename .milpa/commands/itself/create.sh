#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

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

function createOrExit() {
  [[ -f "$1" ]] && @milpa.fail "Command already exists at $1"

  @milpa.log info "Creating command $(@milpa.fmt bold "${command_name}") at $1"
  mkdir -p "$(dirname "$1")"
}

case "${MILPA_OPT_KIND}" in
  bash)
    spec="${command_path}.yaml"
    command_path="$command_path.sh"
    createOrExit "$command_path"
    echo "#!/usr/bin/env bash" >> "$command_path" || @milpa.fail "could not create $command_path"
    ;;
  zsh)
    spec="${command_path}.yaml"
    command_path="$command_path.zsh"
    createOrExit "$command_path"
    echo "#!/usr/bin/env zsh" >> "$command_path" || @milpa.fail "could not create $command_path"
    ;;
  fish)
    spec="${command_path}.yaml"
    command_path="$command_path.fish"
    createOrExit "$command_path"
    echo "#!/usr/bin/env fish" >> "$command_path" || @milpa.fail "could not create $command_path"
    ;;
  executable)
    createOrExit "$command_path"
    touch "$command_path" || @milpa.fail "could not create $command_path"
    chmod +x "$command_path"
    spec="${command_path}.yaml"
    ;;
esac

# shellcheck disable=2001
cat > "$spec" <<YAML
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
