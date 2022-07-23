#!/bin/env bash

declare -a _tmp_files_created

function @tmp.file () {
  local prefix; prefix="${1}"
  fname=$(mktemp "/tmp/$prefix.XXXXXX")
  _files_created=( "$fname" )
  printf -v "$1" -- '%s' "$fname"
}

function @tmp.dir () {
  local prefix; prefix="${1}"
  fname=$(mktemp -d "/tmp/$prefix.XXXXXX")
  _files_created+=( "$fname" )
  printf -v "$1" -- '%s' "$fname"
}

function @tmp.cleanup () {
  @milpa.log debug "cleaning up created files: ${_files_created[*]}"
  for file in "${_files_created[@]}"; do
    if [[ -f "$file" ]]; then
      @milpa.log debug "removing temporary file $file"
      rm -f "$file"
    elif [[ -d "$file" ]]; then
      @milpa.log debug "removing temporary dir $file"
      rm -rf "$file"
    fi
  done
}
