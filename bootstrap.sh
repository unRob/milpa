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
# curl -LO https://milpa.dev/install.sh | bash -

if [[ -t 1 ]] && [[ -z ${NO_COLOR+x} ]]; then
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

META_BASE="${META_BASE:-https://milpa.dev/}/.well-known/milpa/"
ASSET_BASE="${GITHUB_REPO:-"https://github.com/unRob/milpa"}/releases" #/latest/download/ASSET.ext
if [[ -x ${VERSION+x} ]]; then
  >&2 echo "${_FMT_GRAY}No VERSION provided, querying for default${_FMT_RESET}"
  VERSION=$(curl -L "$META_BASE/latest-version")
fi
PREFIX="${PREFIX:-/usr/local/lib}"
TARGET="${PREFIX:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin) OS="darwin";;
  Linux) OS="linux";;
  *)
    >&2 echo "${_FMT_ERROR}unsupported OS: $OS${_FMT_RESET}"
    exit 2
esac

machine="$(uname -m)"
case "$machine" in
  x86_64) ARCH="amd64" ;;
  armv7l) ARCH="arm" ;;
  *) ARCH="$machine"
esac

sudo mkdir -pv "$PREFIX"

>&2 echo "${_FMT_BOLD}Downloading milpa version $VERSION to $PREFIX/milpa${_FMT_RESET}"
curl -LO "$ASSET_BASE/$VERSION/dowload/milpa-$OS-$ARCH.tgz" | sudo tar xfz -C "$PREFIX" -
>&2 echo "Installing symbolic links to $TARGET"
sudo ln -sfv "$PREFIX/milpa/milpa" "$TARGET/milpa"
sudo ln -sfv "$PREFIX/milpa/compa" "$TARGET/compa"
sudo mkdir -pv "${PREFIX}/milpa/repos"
mkdir -pv "${XDG_HOME_DATA:-$HOME/.local/share}/milpa/repos"

installed_version=$("$TARGET/milpa" --version) || {
  >&2 echo "${_FMT_ERROR}Could not get the installed version${_FMT_RESET}"
  exit 2
}

header="🌽 Installed milpa version $installed_version 🌽"
hlen="$(( ${#header} + 3 ))"
line="$(printf -- "-%.0s" $(seq 1 "$hlen"))"
>&2 echo "$line"
>&2 echo "${_FMT_INVERTED}$header$_FMT_RESET"
>&2 echo "$line"
>&2 echo "${_FMT_WARNING}Run 'milpa itself shell install-autocomplete' to install shell completions${_FMT_RESET}"
