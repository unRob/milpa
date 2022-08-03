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

# Run with
# curl -L https://milpa.dev/install.sh | bash -

if [[ -t 1 ]] && [[ -z ${NO_COLOR+x} ]]; then
  [[ -z ${TERM+x} ]] && export TERM="${TERM:-xterm-color}"
  _FMT_INVERTED="$(tput rev)"
  _FMT_BOLD="$(tput bold)"
  _FMT_RESET="$(tput sgr0)"
  _FMT_ERROR="$(tput setaf 1)"
  _FMT_WARNING="$(tput setaf 3)"
  _FMT_GRAY="$(tput setaf 7)"
else
  _FMT_INVERTED=""
  _FMT_BOLD=""
  _FMT_RESET=""
  _FMT_ERROR=""
  _FMT_WARNING=""
  _FMT_GRAY=""
fi

function @fail () {
  set +o xtrace
  # print an error, then exit
  >&2 echo "${_FMT_ERROR}$*${_FMT_RESET}"
  exit 2
}

function @info () {
  >&2 echo "${@}"
}

# a place to look for a well-known latest version
MILPA_META_BASE="${MILPA_META_BASE:-https://milpa.dev/}/.well-known/milpa/"
# a github repo to pull assets from
ASSET_BASE="${MILPA_GITHUB_REPO:-"https://github.com/unRob/milpa"}/releases" #/latest/download/ASSET.ext
# version can be set by specifying MILPA_VERSION, otherwise we'll find out from the internet
if [[ "${MILPA_VERSION}" == "" ]]; then
  >&2 echo "${_FMT_GRAY}No VERSION provided, querying for default${_FMT_RESET}"
  MILPA_VERSION=$(curl --silent --fail --show-error -L "$MILPA_META_BASE/latest-version") || @fail "Could not fetch latest version!"
fi
# Where the package gets installed to
PREFIX="${PREFIX:-/usr/local/lib}/milpa"
# Where we drop links to binaries at
TARGET="${TARGET:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin) OS="darwin";;
  Linux) OS="linux";;
  *) @fail "No builds available for $OS, only darwin and linux"
esac

machine="$(uname -m)"
case "$machine" in
  x86_64) ARCH="amd64" ;;
  armv7l) ARCH="arm" ;;
  aarch64) ARCH="arm64" ;;
  *) ARCH="$machine"
esac

case "$ARCH" in
  amd64|arm|arm64|mips) @info "Detected system: $OS/$ARCH";;
  *) @fail "No builds available for $OS/$ARCH"
esac

# where system-level repos live
globalRepos="${PREFIX}/repos"
default_data_home="$HOME/.local/share"
# user-specific milpa-related files go here
milpaLocal="${XDG_DATA_HOME:-$default_data_home}/milpa"
localRepos="${milpaLocal}/repos"
package="milpa-$OS-$ARCH.tgz"

# Get the package
if [[ ! -f "$package" ]]; then
  @info "${_FMT_BOLD}Downloading milpa version $VERSION to $PREFIX${_FMT_RESET}"
  curl --silent --fail --show-error -LO "$ASSET_BASE/download/$VERSION/$package" || @fail "Could not download milpa package from $ASSET_BASE/download/$VERSION/$package"
else
  @info "${_FMT_BOLD}Using downloaded package at $package${_FMT_RESET}"
fi

@info "Downloaded $ASSET_BASE/download/$VERSION/$package"

# Find some nice spot in the ground
if [[ ! -d "$PREFIX" ]]; then
  @info "${_FMT_BOLD}Creating $PREFIX, enter your password if prompted${_FMT_RESET}"
  if [[ -w "$(dirname "$PREFIX")" ]]; then
    mkdir -pv "$PREFIX"
  else
    sudo mkdir -pv "$PREFIX"
  fi || @fail "Could not create $PREFIX directory"
else
  @info "${_FMT_WARNING}$PREFIX already exists, deleting previous installation...${_FMT_RESET}"
  if [[ -w "$PREFIX" ]]; then
    find "$PREFIX" -maxdepth 1 -mindepth 1 \! -name repos -exec rm -rf {} \;
  else
    sudo find "$PREFIX" -maxdepth 1 -mindepth 1 \! -name repos -exec rm -rf {} \;
  fi
fi

# dig a hole, pour some seeds
if [[ -w "$PREFIX" ]]; then
  tar xfz "$package" -C "$(dirname "$PREFIX")" || @fail "Could not extract milpa package to $PREFIX"
else
  sudo tar xfz "$package" -C "$(dirname "$PREFIX")" || @fail "Could not extract milpa package to $PREFIX"
fi

# get ready for growing some scripts
@info "Installing symbolic links to $TARGET"
if [[ -w "$PREFIX" ]]; then
  ln -sfv "$PREFIX/milpa" "$TARGET/milpa"
  ln -sfv "$PREFIX/compa" "$TARGET/compa"
  [[ -d "$globalRepos" ]] || mkdir -pv "$globalRepos"
else
  sudo ln -sfv "$PREFIX/milpa" "$TARGET/milpa"
  sudo ln -sfv "$PREFIX/compa" "$TARGET/compa"
[[ -d "$globalRepos" ]] || sudo mkdir -pv "$globalRepos"
fi

# recycle the bag
rm -rf "$package"

[[ -d "$localRepos" ]] || mkdir -pv "$localRepos"
# update version so milpa doesn't look for updates innecessarily
date "+%s" > "$milpaLocal/last-update-check"

# Test we can run milpa
installed_version="$("$TARGET/milpa" --version)" || @fail "Could not get the installed version"

header="🌽 Installed milpa version $installed_version 🌽"
hlen="$(( ${#header} + 3 ))"
line="$(printf -- "-%.0s" $(seq 1 "$hlen"))"
@info "$line"
@info "${_FMT_INVERTED}$header$_FMT_RESET"
@info "$line"
@info "${_FMT_WARNING}Run 'milpa itself shell install-autocomplete' to install shell completions${_FMT_RESET}"
