#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

export MILPA_VERSION="$MILPA_ARG_VERSION"

output="${MILPA_OPT_OUTPUT:-$MILPA_ROOT/dist}"
all_targets=( linux/amd64 linux/arm64 linux/arm linux/mips darwin/amd64 darwin/arm64 )
# build packages

if [[ "${#MILPA_ARG_TARGETS}" -eq 0 ]] || [[ "${MILPA_ARG_TARGETS[1]}" == "auto" ]]; then
  MILPA_ARG_TARGETS=( "${all_targets[@]}" )
fi

@milpa.log info "Starting build for ${MILPA_ARG_TARGETS[*]}"
mkdir -p "$output"
CGO_ENABLED=0 GOFLAGS="-trimpath" gox -osarch "${MILPA_ARG_TARGETS[*]}" \
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

  cp -r ./milpa ./.milpa LICENSE.txt README.md CHANGELOG.md "$dist_dir/"
  rm -rf "$package"
  tar -czf "$package" -C "$(dirname "$dist_dir")" milpa || @milpa.fail "Could not archive $package"
  openssl dgst -sha256 "$package" | awk '{print $2}' > "${package##.tgz}.shasum" || @milpa.fail "Could not generate shasum for $package"
done
@milpa.log success "Archives generated"


function fetchDoc () {
  curl --silent --fail --show-error "http://localhost:4242$path" || @milpa.fail "Could not fetch $path"
}

function renderDoc() {
  @milpa.log info "Serializing docs for $path"
  mkdir -p "${html}${1:-}"

  set -o pipefail
  fetchDoc "$path" |
    tidy -quiet -wrap 0 -indent - > "${html}/${path}${path+/}index.html"

  [[ "$?" -gt 1 ]] && @milpa.fail "Could not tidy up $path"
  return 0
}

@milpa.log info "Generating html docs"
mp="$MILPA_PATH"
html="$output/${MILPA_ARG_HOSTNAME##*//}"
export MILPA_DISABLE_USER_REPOS=true
export MILPA_DISABLE_GLOBAL_REPOS=true
cat - <(tail -n +2 "$MILPA_ROOT/CHANGELOG.md") > "$MILPA_ROOT/.milpa/docs/milpa/changelog.md" <<YAML
---
description: "Changelog entries for every released version"
weight: 100
---

YAML
MILPA_PATH="" "$MILPA_COMPA" help docs --server --base "$MILPA_ARG_HOSTNAME" &
pid=$!
@milpa.log info "started server at pid $pid"
trap 'rm "$MILPA_ROOT/.milpa/docs/milpa/changelog.md"; kill -9 "$pid"' ERR EXIT
sleep 3

mkdir "$html"
renderDoc "/"

while read -r path; do
  renderDoc "$path" || @milpa.fail "Could not fetch $path"
done < <(htmlq --attribute href "#commands a" <dist/milpa.dev/index.html)

# this returns a 404, which is very much expected so no --show-error nor --fail
curl --silent "http://localhost:4242/404" | tidy -quiet -wrap 0 -indent - > "${html}/404.html"

unset MILPA_DISABLE_USER_REPOS MILPA_DISABLE_GLOBAL_REPOS
export MILPA_PATH="$mp"
@milpa.log success "Docs exported"

@milpa.log info "Copying website assets"
# static files
cp -r internal/docs/static "$html/static"
# Copy over bootstrap script
cp "$MILPA_ROOT/bootstrap.sh" "$html/install.sh"
# Write version to a well-known location
mkdir -p "$html/.well-known/milpa"
echo -n "$MILPA_VERSION" > "$html/.well-known/milpa/latest-version"
# github pages needs a CNAME, provide one
echo -n "${MILPA_ARG_HOSTNAME##*//}" > "$html/CNAME"
# github pages doesn't need to process our docs as jekyll
echo -n "${MILPA_ARG_HOSTNAME}" > "$html/.nojekyll"
@milpa.log success "HTML docs written to $html"

@milpa.log complete "Release built to $output"
