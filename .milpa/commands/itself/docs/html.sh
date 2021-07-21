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

@milpa.log info "generating command docs"
MILPA_PLAIN_HELP=enabled "$MILPA_COMPA" __generate_documentation "$content" || @milpa.fail "Could not generate command documentation"

mv -v "$content/help/docs.md" "$content/help/docs/_index.md"

cat - <(tail -n +2 "$MILPA_ROOT/CHANGELOG.md") > "$content/help/docs/milpa/changelog.md" <<YAML
---
weight: 100
---

YAML


if [[ "$MILPA_ARG_ACTION" == "serve" ]]; then
  containerID="$(docker ps -q --filter name=milpa_docs)"
  if [[ "$containerID" != "" ]]; then
    @milpa.log info "Website generator already up at <$containerID>"
    exec docker attach "$containerID"
  fi

  @milpa.log info "Launching hugo website generator with $MILPA_OPT_IMAGE"
  containerID="$(docker run --rm -it \
    --name milpa_docs \
    --detach \
    -p "$MILPA_OPT_PORT:$MILPA_OPT_PORT" \
    -v "$content:/src/content/" \
    "$MILPA_OPT_IMAGE" serve --debug --port "$MILPA_OPT_PORT")" || @milpa.fail "Could not spin up website generator"
    # to debug add:
    # -v "${MILPA_ROOT}/repos/internal/docs/.template/theme:/src/themes/cli/" \
    # -v "${MILPA_ROOT}/repos/internal/docs/.template/config.toml:/src/config.toml" \

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
    -v "${dst}:/src/public" \
    "$MILPA_OPT_IMAGE"
fi
