#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

tmpdir="${DOCS_TMP_DIR:-$HOME/.cache/milpa}"
content="$tmpdir/html"

function generate_content_folder() {
  @milpa.log info "generating command docs"
  MILPA_PLAIN_HELP=enabled "$MILPA_COMPA" __generate_documentation "$content" || @milpa.fail "Could not generate command documentation"

  mv -v "$content/help/docs.md" "$content/help/docs/_index.md"

  mkdir -p "$content/help/docs/milpa"
  cat - <(tail -n +2 "$MILPA_ROOT/CHANGELOG.md") > "$content/help/docs/milpa/changelog.md" <<YAML
---
description: "Changelog entries for every released version"
weight: 100
---

YAML
}

generate_content_folder


if [[ "$MILPA_ARG_ACTION" == "serve" ]]; then
  containerID="$(docker ps -q --filter name=milpa_docs)"
  if [[ "$containerID" != "" ]]; then
    @milpa.log info "Website generator already up at <$containerID>"
    exec docker attach "$containerID"
  fi

  @milpa.log info "Launching hugo website generator with $MILPA_OPT_IMAGE"

  templateDir=()
  if [[ "$MILPA_OPT_TEMPLATE" != "" ]]; then
    tplDir="$(cd "$MILPA_OPT_TEMPLATE" && pwd)" || @milpa.fail "--template is not a directory"
    templateDir=(
      -v "$tplDir/theme:/src/themes/cli/"
      -v "$tplDir/config.toml:/src/config.toml"
    )
  fi

  containerID="$(docker run --rm -it \
    --name milpa_docs \
    --detach \
    -p "$MILPA_OPT_PORT:$MILPA_OPT_PORT" \
    -v "$content:/src/content/" \
    "${templateDir[@]}" \
    "$MILPA_OPT_IMAGE" serve --debug --port "$MILPA_OPT_PORT")" || @milpa.fail "Could not spin up website generator"
    # to debug add:

  docker logs --follow milpa_docs &

  @milpa.log info "Waiting for server to be up"
  while ! curl -s --fail "http://localhost:$MILPA_OPT_PORT" >/dev/null 2>&1; do
    sleep 1
  done
  @milpa.log complete 'Server is up, attaching container and opening address'

  [[ "$OSTYPE" == darwin* ]] && open "http://localhost:$MILPA_OPT_PORT"
  [[ "$OSTYPE" == linux* ]] && xdg-open "http://localhost:$MILPA_OPT_PORT"

  (
    trap 'kill 0' SIGINT;
    if command -v fswatch >/dev/null; then
      IFS=: read -ra MILPA_PATH_ARR <<<"$MILPA_PATH"
      @milpa.log info "Listening for changes in ${MILPA_PATH_ARR[*]//:/ }..."
      @milpa.log warning "Press CTRL-C twice to stop"
      fswatch --one-per-batch --recursive --print0 "${MILPA_PATH_ARR[@]//:/ }" | while read -r -d "" _; do
        @milpa.log info "Changes found on MILPA_PATH"
        generate_content_folder
      done &
    else
      @milpa.log warning "fswatch is not available, will not listen for changes"
    fi
    docker attach "$containerID"
    wait
  )
else
  dst="$(realpath "${MILPA_OPT_TO}")/${MILPA_OPT_HOSTNAME}"
  @milpa.log info "Writing docs to $dst"
  mkdir -p "$dst"
  docker run --rm \
    -v "$content:/src/content/" \
    -v "${dst}:/src/public" \
    "$MILPA_OPT_IMAGE" || @milpa.fail "Could not write docs"
fi
