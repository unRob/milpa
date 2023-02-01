#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

cd "$MILPA_ROOT" || @milpa.fail "could not cd into $MILPA_ROOT"

@milpa.log info "Linting shell scripts"
find .milpa repos/internal/commands -name '*.sh' -exec shellcheck {} \+  || @milpa.fail "could not lint commands"
shellcheck bootstrap.sh test/_helpers/milpa/*.bash || @milpa.fail "could not lint helper files"
@milpa.log complete "Shell files are up to spec"
