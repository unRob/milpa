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
