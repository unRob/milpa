#!/usr/bin/env bash

if [[ "$MILPA_OPT_REPO" == "" ]]; then
  repo_path="$(pwd)"
  while [[ ! -d "$repo_path/.milpa" ]]; do
    repo_path=$(dirname "$repo_path")
  done
else
  repo_path="$MILPA_OPT_REPO"
fi

milpa="$repo_path/.milpa"

IFS="/"  path="$milpa/commands/$(echo "${MILPA_ARG_NAME[*]}")"
mkdir -p "$(dirname "$path")"

if [[ "${MILPA_OPT_EXECUTABLE}" ]]; then
  touch "$path"
  chmod +x "$path"
else
  echo "#!/usr/bin/env bash" >> "$path.sh"
fi

cat > "$path.${MILPA_OPT_CONFIG_FORMAT}" <<YAML
summary: Does a thing
description: |
  Longer description of how it does the thing
# arguments:
#   - name: something
#     description: passes something to your script
#     set:
#       from: { subcommand: another command }
# options:
#   explode:
#     type: boolean
#     description: something else
YAML
