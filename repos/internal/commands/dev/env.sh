#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

@milpa.load_util shell
@milpa.log info "Exporting shell variables"

root="${MILPA_COMMAND_REPO%%/repos/internal*}"
@milpa.shell.export "MILPA_ROOT" "$root"

# shellcheck disable=2049
if [[ ! "$PATH" = "$root"* ]]; then
  @milpa.shell.prepend_path "$root"
fi
