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

cores="$(sysctl -n hw.physicalcpu 2>/dev/null || grep -c ^processor /proc/cpuinfo)"

export MILPA_VERSION="$MILPA_ARG_VERSION"
@milpa.log info "Starting build"

output="${MILPA_OPT_OUTPUT:-$MILPA_ROOT/dist}"

# build packages
RELEASE_TARGET="$output" make -j"$cores" "$output/packages" || @milpa.fail "Could not complete release build"
@milpa.log success "Build complete"

# create docs
milpa release docs-image --skip-publish "$MILPA_VERSION" || @milpa.fail "Could not build docs image"
@milpa.log info "Generating html docs"
MILPA_PATH="$MILPA_ROOT/.milpa" milpa itself docs html write \
  --to "$output" \
  --image milpa-docs \
  --hostname "$MILPA_ARG_HOSTNAME" || @milpa.fail "Could not generate docs"
@milpa.log success "Docs exported"

@milpa.log info "Copying website assets"
html="$output/$MILPA_ARG_HOSTNAME"
# Copy over bootstrap script
cp "$MILPA_ROOT/bootstrap.sh" "$html/install.sh"
# Write version to a well-known location
mkdir -p "$html/.well-known/milpa"
echo -n "$MILPA_VERSION" > "$html/.well-known/milpa/latest-version"
# github pages needs a CNAME, provide one
echo -n "$MILPA_ARG_HOSTNAME" > "$html/CNAME"
@milpa.log success "HTML docs written to $html"

@milpa.log complete "Release built to $output"
