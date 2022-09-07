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

export MILPA_VERSION="$MILPA_ARG_VERSION"

output="${MILPA_OPT_OUTPUT:-$MILPA_ROOT/dist}"
all_targets=( linux/amd64 linux/arm64 linux/arm linux/mips darwin/amd64 darwin/arm64 )
# build packages

if [[ "${#MILPA_ARG_TARGETS}" -eq 0 ]] || [[ "${MILPA_ARG_TARGETS[1]}" == "auto" ]]; then
  MILPA_ARG_TARGETS=( "${all_targets[@]}" )
fi

@milpa.log info "Starting build for ${MILPA_ARG_TARGETS[*]}"
mkdir -p "$output"
GOFLAGS="-trimpath" gox -osarch "${MILPA_ARG_TARGETS[*]}" \
  -parallel="$MILPA_OPT_PARALLEL" \
  -ldflags "-s -w -X main.version=${MILPA_VERSION}" \
  -output "$output/{{.OS}}-{{.Arch}}" || @milpa.fail "Could not build with gox"
@milpa.log success "Build complete"


@milpa.log info "Generating archives"
for pair in "${MILPA_ARG_TARGETS[@]//\//-}"; do
  dist_dir="$output/tmp/$pair/milpa"
  package="$output/milpa-$pair.tgz"

  mkdir -p "$dist_dir"
  if [[ ! -f "$dist_dir/compa" ]]; then
    if [[ "$pair" != "darwin-arm64" ]]; then
      upx --no-progress --best -o "$dist_dir/compa" "$output/$pair" || @milpa.fail "Could not compress $dist_dir/compa"
    else
      @milpa.warning "UPX produces botched arm64 builds :/"
      @milpa.warning https://github.com/upx/upx/issues/446
      cp "$output/$pair" "$dist_dir/compa"
    fi
    rm -rf "${output:?}/$pair"
  fi

  cp -rv ./milpa ./.milpa LICENSE.txt README.md CHANGELOG.md "$dist_dir/"
  rm -rf "$package"
  tar -czf "$package" -C "$(dirname "$dist_dir")" milpa || @milpa.fail "Could not archive $package"
  openssl dgst -sha256 "$package" | awk '{print $2}' > "${package##.tgz}.shasum" || @milpa.fail "Could not generate shasum for $package"
done
@milpa.log success "Archives generated"

# create docs
milpa release docs-image --skip-publish "$MILPA_VERSION" || @milpa.fail "Could not build docs image"
@milpa.log info "Generating html docs"
mp="$MILPA_PATH"
export MILPA_DISABLE_USER_REPOS=true
export MILPA_DISABLE_GLOBAL_REPOS=true
MILPA_PATH="" milpa itself docs html write \
  --to "$output" \
  --image milpa-docs \
  --hostname "$MILPA_ARG_HOSTNAME" || @milpa.fail "Could not generate docs"
unset MILPA_DISABLE_USER_REPOS MILPA_DISABLE_GLOBAL_REPOS
export MILPA_PATH="$mp"
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
# github pages doesn't need to process our docs as jekyll
echo -n "$MILPA_ARG_HOSTNAME" > "$html/.nojekyll"
@milpa.log success "HTML docs written to $html"

@milpa.log complete "Release built to $output"
