#!/usr/bin/env bash
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

base="$(git rev-parse --show-toplevel)"
cd "$base" || @milpa.fail "could not cd into root directory"

@milpa.log info "Configuring git hooks"
git config core.hooksPath "$base/internal/bin/hooks"
@milpa.log info "Fetching notes"
git fetch origin refs/notes/*:refs/notes/*

@milpa.log info "Making sure submodules are here"
git submodule update --init --recursive

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
fi

@milpa.log info "Installing go packages"
packages=(
  gotest.tools/gotestsum@v1.8.1
  github.com/mitchellh/gox@v1.0.1
  github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.3
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

exec milpa dev build
