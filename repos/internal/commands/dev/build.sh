#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

if [[ "$MILPA_ARG_VERSION" == "auto" ]]; then
  VERSION="$(git describe --long)"
else
  VERSION="$MILPA_ARG_VERSION"
fi
cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"

args=(
  -ldflags "-X main.version=${VERSION}" -o milpa
)

if [[ "${MILPA_OPT_COVERAGE}" ]]; then
  # args+=(-cover "-coverpkg=$(go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' -deps . |
  #   grep -E '(milpa|chinampa)' |
  #   paste -sd "," -)" )
  @milpa.log info "Collecting coverage"
  args+=( -cover -coverpkg=./... -tags coverage )
fi

@milpa.log info "Building milpa version $VERSION"
go build "${args[@]}" || exit 2
# account older milpa versions depending on bash entrypoint
ln -sfv milpa compa
@milpa.log complete "milpa version $VERSION built"

