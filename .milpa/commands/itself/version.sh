#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

"$MILPA" __version 2>&1 || {
  if [[ "$?" != 42 ]]; then
    @milpa.fail "could not get version"
  fi
}

if [[ -t 1 ]]; then
  echo ""
fi
