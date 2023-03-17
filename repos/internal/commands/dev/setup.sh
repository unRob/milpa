#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

base="$(git rev-parse --show-toplevel)"
cd "$base" || @milpa.fail "could not cd into root directory"

@milpa.log info "Configuring git hooks"
git config core.hooksPath "$base/internal/bin/hooks"
@milpa.log info "Fetching notes"
git fetch origin refs/notes/*:refs/notes/*

@milpa.log info "Making sure submodules are here"
git submodule update --init --recursive


## golang
if [[ "$ASDF_DIR" ]]; then
  @milpa.log info "Installing golang version with asdf"

  if ! asdf plugin list | grep golang >/dev/null; then
    asdf plugin add golang || @milpa.fail "Could not install golang plugin"
  fi

  if ! asdf list golang | grep -f <(cut -d" " -f 2 .tool-versions) >/dev/null; then
    asdf install || @milpa.fail "could not install golang version"
    asdf reshim golang
  fi
  @milpa.log success "go is now installed"
else
  if command -v golang >/dev/null; then
    installed=$(go version | awk '{gsub(/[a-z]/, "", $3); print($3)}' 2>/dev/null)
    required=$(awk '/golang/ {print $2}' .tool-versions)
    if [[ "$installed" != "$required" ]]; then
      @milpa.fail "Golang v${required} is not installed (found v${installed}), please install it from https://go.dev/doc/install"
    fi

    @milpa.log success "Go v${required} is already installed"
  else
    @milpa.fail "Golang is not installed, please install go v${required}: https://go.dev/doc/install"
  fi
fi

@milpa.log info "Installing go packages"
packages=(
  gotest.tools/gotestsum@v1.9.0
  github.com/mitchellh/gox@9f71238 # used to be 1.0.1 but there doesn't seem to be more releases?
  github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2
)

for package in "${packages[@]}"; do
  name="$(basename "$package")"
  bin="${name##@*}"
  if command -v "$bin" >/dev/null; then
    @milpa.log success "$package already installed"
    continue
  fi

  @milpa.log info "Installing $package"
  go install "$package" || @milpa.fail "Could not install $package"
  @milpa.log success "Installed $package"
done

[[ -d "$ASDF_DIR" ]] && asdf reshim golang

go mod tidy || @milpa.fail "go mod tidy failed"

## shellcheck
if ! command -v shellcheck >/dev/null; then
  @milpa.fail "Shellcheck is not installed, see https://github.com/koalaman/shellcheck#installing"
fi

exec milpa dev build
