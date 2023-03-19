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
  -ldflags "-s -w -X main.version=${VERSION}" -o compa
)

if [[ "${MILPA_OPT_COVERAGE}" ]]; then
  args+=(-cover -coverpkg=./...)
fi

@milpa.log info "Building compa version $VERSION"
go build "${args[@]}" || exit 2
@milpa.log complete "compa version $VERSION built"

