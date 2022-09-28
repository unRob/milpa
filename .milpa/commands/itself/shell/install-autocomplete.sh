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

caveats=""
case "$SHELL" in
  *bash)
    @milpa.log info "Bash detected"
    if [[ -d /etc/bash_completion.d ]]; then
      "$MILPA_COMPA" __generate_completions bash > "/etc/bash_completion.d/milpa" || @milpa.fail "Could not install completions"
    elif [[ -d /usr/local/etc/bash_completion.d ]]; then
      "$MILPA_COMPA" __generate_completions bash | sudo tee "/usr/local/etc/bash_completion.d/milpa" || @milpa.fail "Could not install completions"
    else
      @milpa.fail "No directory found for writing completion script (tried /etc/bash_completion.d and /usr/local/etc/bash_completion.d)"
    fi
    ;;
  *zsh)
    @milpa.log info "ZSH detected"
    if $SHELL -i -c "command -v compinit >/dev/null" >/dev/null 2>&1; then
      caveats='compinit has not been loaded into this shell, enable it by running:

echo "autoload -U compinit; compinit" >> ~/.zshrc && source ~/.zshrc

then reloading your shell'
    fi

    # shellcheck disable=2016
    dst=$(zsh -i -c 'printf "%s" "${${fpath[@]:#$HOME/*}[1]}"' 2>/dev/null) || @milpa.fail "Unable to locate an fpath to install completions to"
    if [[ -w "$dst" ]]; then
      [[ -f "$dst" ]] || mkdir -pv "$dst"
      "$MILPA_COMPA" __generate_completions zsh > "${dst}/_milpa" || @milpa.fail "Could not install completions"
    else
      @milpa.log warning "$dst does not look writeable for $USER, using sudo"
      [[ -f "$dst" ]] || sudo mkdir -pv "$dst"
      "$MILPA_COMPA" __generate_completions zsh | sudo tee "${dst}/_milpa" >/dev/null || @milpa.fail "Could not install completions"
    fi
    ;;
  *fish)
    @milpa.log info "Fish detected"
    "$MILPA_COMPA" __generate_completions fish > "$HOME/.config/fish/completions/milpa.fish" || @milpa.fail "Could not install completions"
  ;;
  *)
    @milpa.fail "No completion script found for shell $SHELL"
esac

@milpa.log complete "Shell completion added for $SHELL successfully"
[[ "$caveats" == "" ]] || @milpa.log warning "$caveats"
exit 0
