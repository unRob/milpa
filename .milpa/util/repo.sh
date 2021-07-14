#!/usr/bin/env bash

function @milpa.repo.current_path () {
  repo_path="$(pwd)"
  while [[ ! -d "$repo_path/.milpa" ]]; do
    repo_path=$(dirname "$repo_path")
  done
  echo "$repo_path"
}
