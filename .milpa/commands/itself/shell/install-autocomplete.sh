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

case "$SHELL" in
  *bash)
    if [[ -d /etc/bash_completion.d ]]; then
      set -x
      "$MILPA_COMPA" __generate_completions bash > "/etc/bash_completion.d/$MILPA_NAME"
      set +x
    elif [[ -d /usr/local/etc/bash_completion.d ]]; then
      set -x
      "$MILPA_NAME_COMPA" __generate_completions bash > "/usr/local/etc/bash_completion.d/milpa"
      set +x
    else
      @milpa.fail "No directory found for writing completion script (tried /etc/bash_completion.d and /usr/local/etc/bash_completion.d)"
    fi
    ;;
  *zsh)
    $SHELL -i -c "command -v compinit >/dev/null" >/dev/null || @milpa.log warning <<EOF
compinit has not been loaded into this shell, enable it by running

echo "autoload -U compinit; compinit" >> ~/.zshrc
and reloading your shell
EOF
    # shellcheck disable=2016
    dst=$(zsh -i -c 'printf "%s" "${${fpath[@]:#$HOME/*}[1]}"') || @milpa.fail "Unable to locate an fpath to install completions to"
    if [[ -w "$dst" ]]; then
      set -ex
      "$MILPA_COMPA" __generate_completions zsh > "${dst}/_${MILPA_NAME}"
      set +ex
    else
      @milpa.log warning "$dst does not look writeable for $USER, using sudo"
      set -ex
      "$MILPA_COMPA" __generate_completions zsh | sudo tee "${dst}/_${MILPA_NAME}" >/dev/null
      set +ex
    fi

    @milpa.log warning "Please restart your shell"
    ;;
  *fish)
    set -ex
    "$MILPA_COMPA" __generate_completions fish > "$HOME/.config/fish/completions/${MILPA_NAME}.fish"
    set +x
  ;;
  *)
    @milpa.fail "No completion script found for shell $SHELL"
esac

@milpa.log complete "Shell completion added for $SHELL successfully"
