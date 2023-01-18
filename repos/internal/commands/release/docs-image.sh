#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

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
  repo="${MILPA_OPT_DOCKER_REPO%%/*}"
  docker login "$repo" -u "$username" --password-stdin <<<"${!pass_var}" || @milpa.fail "Could not login to the <${repo}> docker repository using username: <${username}> and password from env var: <${pass_var}>"
fi

@milpa.log info "Publishing image to $MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION"
docker image tag milpa-docs "$MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION" || @milpa.fail "Failed tagging image"
docker push "$MILPA_OPT_DOCKER_REPO:$MILPA_ARG_VERSION" || @milpa.fail "Failed pushing image"

if [[ ! "$MILPA_OPT_SKIP_LATEST" ]]; then
  docker image tag milpa-docs "$MILPA_OPT_DOCKER_REPO:latest" || @milpa.fail "Failed tagging latest image"
  docker push "$MILPA_OPT_DOCKER_REPO:latest" || @milpa.fail "Failed pushing latest image"
fi

@milpa.log complete "Image built and published"
