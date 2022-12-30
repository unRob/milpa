#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"

@milpa.log info "Linting go files"
golangci-lint run || exit 2
@milpa.log complete "Go files are up to spec"

