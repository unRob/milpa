#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

runs=()
while IFS=  read -r -d $'\0'; do
    runs+=("$REPLY")
done < <(find test/coverage -type d -maxdepth 1 -mindepth 1 -print0)
packages="$(IFS=, ; echo "${runs[*]}")"


@milpa.log info "Building coverage report from runs: ${runs[*]}"
go tool covdata textfmt -i="$packages" -o test/coverage.cov || @milpa.fail "could not merge runs"
go tool cover -html=test/coverage.cov -o test/coverage.html || @milpa.fail "could not build reports"

@milpa.log complete "Coverage report built"
go tool covdata percent -i="$packages"
go tool cover -func=test/coverage.cov | tail -n 1
