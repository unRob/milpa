#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

milpa dev lint go && milpa dev test unit
goOk=$?

milpa dev lint shell
shellOk="$?"
integrationOk=1
if [[ "$goOk" -eq 0 ]]; then
  milpa dev build
  milpa dev test integration
  integrationOk="$?"
fi

MILPA_PATH="" MILPA_DISABLE_USER_REPOS="true" milpa itself doctor --summary
doctorOk="$?"

[[ "${shellOk}${goOk}${integrationOk}${doctorOk}" == "0000" ]]
