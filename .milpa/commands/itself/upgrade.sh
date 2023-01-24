#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
@milpa.load_util milpa-version

@milpa.log info "Looking for latest version..."
# on slow connections, a longer time might be needed to perform the check
if versions="$(@milpa.version.is_latest 10)"; then
  @milpa.fail "Already running latest version"
fi

read -r latest installed <<<"$versions"

@milpa.log success "Upgrading to version $latest from $installed"
export MILPA_VERSION="$latest"
exec curl -L https://milpa.dev/install.sh | bash -
