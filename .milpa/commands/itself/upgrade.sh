#!/usr/bin/env bash
@milpa.load_util milpa-version

if versions="$(@milpa.version.is_latest)"; then
  @milpa.fail "Already running latest version"
fi

read -r latest installed <<<"$versions"

@milpa.log success "Upgrading to version $latest from $installed"
export VERSION="$latest"
exec curl -L https://milpa.dev/install.sh | bash -
