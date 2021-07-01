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

META_BASE="${DOWNLOAD_BASE:-https://milpa.dev/release/}"
ASSET_BASE="${GITHUB_REPO:-"https://github.com/unRob/milpa"}/releases" #/latest/download/ASSET.ext
if [[ -x ${VERSION+x} ]]; then
  >&2 echo "No VERSION provided, querying for default"
  VERSION=$(curl -L "$META_BASE/meta/latest-version")
fi
PREFIX="${PREFIX:-/usr/local/lib}"
TARGET="${PREFIX:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin) OS="darwin";;
  Linux) OS="linux";;
  *)
    >&2 echo "unsupported OS: $OS"
    exit 2
esac

machine="$(uname -m)"
case "$machine" in
  x86_64) ARCH="amd64" ;;
  armv7l) ARCH="arm" ;;
  *) ARCH="$machine"
esac

sudo mkdir -pv "$PREFIX"

>&2 echo "Downloading milpa v$VERSION to $PREFIX/milpa"
curl -LO "$ASSET_BASE/$VERSION/dowload/milpa.tgz" | sudo tar xfz -C "$PREFIX" -
>&2 echo "Downloading compa to $PREFIX/milpa/compa"
curl -LO "$ASSET_BASE/$VERSION/dowload/compa-$OS-$ARCH.tgz" | sudo tar xfz -C "$PREFIX/milpa" -
>&2 echo "Installing symbolic link to $TARGET/milpa"
sudo ln -sfv "$PREFIX/milpa/milpa" "$TARGET/milpa"
