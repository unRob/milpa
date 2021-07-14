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

milpa release docs-image \
  --docker-login "$DOCKER_USERNAME:DOCKER_PASSWORD" \
  "$MILPA_ARG_VERSION" || @milpa.fail "Could not release docs docker image"

release_url="https://api.github.com/repos/unRob/milpa/releases/tags/${MILPA_ARG_VERSION}"

set -o pipefail
UPLOAD_URL="$(curl --silent --fail --show-error "$release_url" | jq -r '.upload_url | split("{")[0]')" || @milpa.fail "Could not determine upload URL for tag ${MILPA_ARG_VERSION}"
@milpa.log info "Uploading assets to $UPLOAD_URL"

find "${MILPA_OPT_OUTPUT}/packages" -name '*.tgz' | while read -r package; do
  fname=$(basename "$package")
  label="${fname%.*}"
  platform=${label##*milpa-}

  @milpa.log info "Uploading $platform binary to release"
  curl --silent \
    --show-error \
    --fail \
    -XPOST \
    -H 'Content-type: application/gzip' \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    --data-binary @"$package" \
    "${UPLOAD_URL}?name=$fname&label=$label" >/dev/null || @milpa.fail "Could not upload $package"

  @milpa.log info "Uploading $platform shasum to release"
  curl --silent \
    --show-error \
    --fail \
    -XPOST \
    -H 'Content-type: plain/text' \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    --data-binary @"${package%%.tgz}.shasum" \
    "${UPLOAD_URL}?name=$label.shasum&label=$label.shasum" >/dev/null || @milpa.fail "Could not upload $label shasum"

  @milpa.log success "Uploaded $fname"
done

@milpa.log complete "Release is out"
