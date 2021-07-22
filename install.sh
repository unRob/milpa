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
  [[ -z ${TERM+x} ]] && export TERM="xterm-color"
  _FMT_INVERTED=$(tput rev)
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

META_BASE="${META_BASE:-https://milpa.dev/}/.well-known/milpa/"
ASSET_BASE="${GITHUB_REPO:-"https://github.com/unRob/milpa"}/releases" #/latest/download/ASSET.ext
if [[ "${VERSION}" == "" ]]; then
  >&2 echo "${_FMT_GRAY}No VERSION provided, querying for default${_FMT_RESET}"
  VERSION=$(curl --silent --fail --show-error -L "$META_BASE/latest-version") || @fail "Could not fetch latest version!"
fi
PREFIX="${PREFIX:-/usr/local/lib}/milpa"
TARGET="${TARGET:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin) OS="darwin";;
  Linux) OS="linux";;
  *) @fail "unsupported OS: $OS"
esac

machine="$(uname -m)"
case "$machine" in
  x86_64) ARCH="amd64" ;;
  armv7l) ARCH="arm" ;;
  *) ARCH="$machine"
esac


globalRepos="${PREFIX}/repos"
localRepos="${XDG_HOME_DATA:-$HOME/.local/share}/milpa/repos"
package="milpa-$OS-$ARCH.tgz"

# Get the package
if [[ ! -f "$package" ]]; then
  >&2 echo "${_FMT_BOLD}Downloading milpa version $VERSION to $PREFIX${_FMT_RESET}"
  curl --silent --fail --show-error -LO "$ASSET_BASE/download/$VERSION/$package" || @fail "Could not download milpa package"
else
  >&2 echo "${_FMT_BOLD}Using downloaded package at $package${_FMT_RESET}"
fi

# Find some nice spot in the ground
if [[ ! -d "$PREFIX" ]]; then
  >&2 echo "${_FMT_BOLD}Creating $PREFIX, enter your password if prompted${_FMT_RESET}"
  sudo mkdir -pv "$PREFIX" || @fail "Could not create $PREFIX directory"
else
  >&2 echo "${_FMT_WARNING}$PREFIX already exists, deleting previous installation...${_FMT_RESET}"
  sudo find "$PREFIX" -maxdepth 1 -mindepth 1 \! -name repos -exec rm -rf {} \;
fi

# dig a hole, pour some seeds
sudo tar xfz "$package" -C "$(dirname "$PREFIX")" || @fail "Could not extract milpa package to $PREFIX"

# recycle the bag
rm -rf "$package"

# get ready for growing some scripts
>&2 echo "Installing symbolic links to $TARGET"
sudo ln -sfv "$PREFIX/milpa" "$TARGET/milpa"
sudo ln -sfv "$PREFIX/compa" "$TARGET/compa"
[[ -d "$globalRepos" ]] || sudo mkdir -pv "$globalRepos"
[[ -d "$localRepos" ]] || mkdir -pv "$localRepos"

# Test we can run milpa
installed_version=$("$TARGET/milpa" --version) || @fail "Could not get the installed version"

header="🌽 Installed milpa version $installed_version 🌽"
hlen="$(( ${#header} + 3 ))"
line="$(printf -- "-%.0s" $(seq 1 "$hlen"))"
>&2 echo "$line"
>&2 echo "${_FMT_INVERTED}$header$_FMT_RESET"
>&2 echo "$line"
>&2 echo "${_FMT_WARNING}Run 'milpa itself shell install-autocomplete' to install shell completions${_FMT_RESET}"
