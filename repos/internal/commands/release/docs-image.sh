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

docker build --tag milpa-docs -f \
  "$MILPA_COMMAND_REPO/docs/.template/Dockerfile" \
  "$MILPA_COMMAND_REPO/docs/.template" || @milpa.fail "could not build image"

if [[ "$MILPA_OPT_SKIP_PUBLISH" ]]; then
  @milpa.log complete "Image built, publishing skipped"
  exit
fi

if [[ "$MILPA_OPT_DOCKER_LOGIN" ]]; then
  username="${MILPA_OPT_DOCKER_LOGIN%:*}"
  pass_var="${MILPA_OPT_DOCKER_LOGIN#*:}"
  repo="${MILPA_OPT_REPO%%/*}"
  docker login "$repo" -u "$username" --password-stdin <<<"${!pass_var}" || @milpa.fail "Could not login to the <${repo}> docker repository using username: <${username}> and password from env var: <${pass_var}>"
fi

@milpa.log "Publishing image to $MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION"
docker image tag milpa-docs "$MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION" || @milpa.fail "Failed tagging image"
docker push "$MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION" || @milpa.fail "Failed pushing image"

if [[ ! "$MILPA_OPT_SKIP_LATEST" ]]; then
  docker image tag milpa-docs "$MILPA_OPT_DOCKER_REPO:latest" || @milpa.fail "Failed tagging latest image"
  docker push "$MILPA_OPT_DOCKER_REPO:latest" || @milpa.fail "Failed pushing latest image"
fi

@milpa.log complete "Image built and published"
