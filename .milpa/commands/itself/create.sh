#!/usr/bin/env bash

milpa="$MILPA_COMMAND_PACKAGE/.milpa"

path="$milpa/${MILPA_ARG_COMMAND_NAME// /-}"
mkdir -p "$(dirname "$path")"

if [[ "${MILPA_ARG_EXECUTABLE}" ]]; then
  touch "$path"
else
  echo "#!/usr/bin/env bash" >> "$path.sh"
fi

cat > "$path.${MILPA_ARG_CONFIG_FORMAT}" <<YAML
summary: Does a thing
description: |
  Longer description of how it does the thing
arguments:
  - name: something
    description: passes something to your script
options:
  explode:
    boolean: true
    description: something else
YAML
