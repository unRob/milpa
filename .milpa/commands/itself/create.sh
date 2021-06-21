#!/usr/bin/env bash

package_path="$(pwd)"
while [[ ! -d "$package_path/.milpa" ]]; do
  package_path=$(dirname "$package_path")
done

milpa="$package_path/.milpa"

IFS="/"  path="$milpa/$(echo "${MILPA_ARG_NAME[*]}")"
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
# options:
#   explode:
#     type: boolean
#     description: something else
YAML
