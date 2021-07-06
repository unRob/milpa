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

tmpdir="${DOCS_TMP_DIR:-$HOME/.cache/milpa}"
content="$tmpdir/html"
mkdir -pv "$content/docs"
mkdir -pv "$content/commands"

for path in ${MILPA_PATH_ARR[*]}; do
  src="${path}/docs"
  [[ ! -d "$src" ]] && continue

  @milpa.log info "copying docs from $src"
  cp -vr "$src"/* "$content/docs/"
done
cat - <(tail -n +2 "$MILPA_ROOT/README.md") > "$content/docs/milpa/index.md" <<YAML
---
title: milpa
weight: 1
---

YAML

cat - <(tail -n +2 "$MILPA_ROOT/CHANGELOG.md") > "$content/docs/milpa/changelog.md" <<YAML
---
title: Changelog
weight: 100
---

YAML

find "$content/docs" -name "*.md" | while read -r doc; do
  if [[ "${doc##*/}" == "index.md" ]]; then
    mv "$doc" "${doc%/*}/_index.md";
    doc="${doc%/*}/_index.md"
  fi

  if [[ "$OSTYPE" == darwin* ]]; then
    sed -i '' -E 's|\(/.milpa\/|(/|g; s|/index.md|/|g; s|\.md([\)#])|\1|g; s|!milpa!|'"$MILPA_NAME"'|g' "$doc"
  else
    sed -i -E 's|\(/.milpa\/|(/|g; s|/index.md|/|g; s|\.md([\)#])|\1|g; s|!milpa!|'"$MILPA_NAME"'|g' "$doc"
  fi
done

@milpa.log info "generating command docs"
MILPA_PLAIN_HELP=enabled "$MILPA_COMPA" __generate_documentation "$content/commands" || @milpa.fail "Could not generate command documentation"


if [[ "$MILPA_ARG_ACTION" == "serve" ]]; then
  containerID="$(docker ps --filter name=milpa_docs -q)"
  if [[ "$containerID" != "" ]]; then
    @milpa.log info "Website generator already up at <$containerID>"
    exec docker attach "$containerID"
  fi


  @milpa.log info "Launching hugo website generator"
  containerID="$(docker run --rm -it \
    --name milpa_docs \
    --detach \
    -p "$MILPA_OPT_PORT:$MILPA_OPT_PORT" \
    -v "$content:/src/content/" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/config.toml:/src/config.toml" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/_variables_project.scss:/src/assets/scss/_variables_project.scss" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/_index.html:/src/content/_index.html" \
    "$MILPA_OPT_IMAGE" serve --port "$MILPA_OPT_PORT")" || @milpa.fail "Could not spin up website generator"

  docker logs --follow milpa_docs &

  @milpa.log info "Waiting for server to be up"
  while ! curl -s --fail "http://localhost:$MILPA_OPT_PORT" >/dev/null 2>&1; do
    sleep 1
  done
  @milpa.log complete 'Server is up, attaching container and opening address'

  [[ "$OSTYPE" == darwin* ]] && open "http://localhost:$MILPA_OPT_PORT"
  [[ "$OSTYPE" == linux* ]] && xdg-open "http://localhost:$MILPA_OPT_PORT"

  exec docker attach "$containerID"
else
  dst="$(realpath "${MILPA_OPT_TO}")/${MILPA_OPT_HOSTNAME}"
  @milpa.log info "Writing docs to $dst"
  mkdir -p "$dst"
  exec docker run --rm -it \
    -v "$content:/src/content/" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/config.toml:/src/config.toml" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/_variables_project.scss:/src/assets/scss/_variables_project.scss" \
    -v "${MILPA_ROOT}/.milpa/docs/.template/_index.html:/src/content/_index.html" \
    -v "${dst}:/src/public" \
    "$MILPA_OPT_IMAGE"
fi
