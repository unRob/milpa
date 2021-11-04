#!/usr/bin/env bash

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

[[ "${shellOk}${goOk}${integrationOk}" == "000" ]]
