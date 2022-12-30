#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

if [[ -z ${!MILPA_ARG_VAR+x} ]]; then
  echo "${MILPA_ARG_VAR} is not set"
  exit 2
fi

echo "${!MILPA_ARG_VAR}"
