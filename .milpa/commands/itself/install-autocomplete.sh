#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

caveats=""
case "$SHELL" in
  *bash)
    @milpa.log info "Bash detected"
    places=(
      /etc/bash_completion.d
      /usr/local/etc/bash_completion.d
    )

    installed=""
    for dst in "${places[@]}"; do
      if [[ -d "$dst" ]]; then
        @milpa.log info "Found completion dir at $dst"
        if [[ -w "$dst" ]]; then
          "$MILPA_COMPA" __generate_completions bash > "$dst/milpa" || @milpa.fail "Could not install completion script"
        else
          @milpa.log warning "$dst does not look writeable for user $USER, using sudo"
          "$MILPA_COMPA" __generate_completions bash | sudo tee "$dst/milpa" || @milpa.fail "Could not install completion script"
        fi
        @milpa.log success "Installed completion script to $dst/milpa"
        installed="true"
        break
      fi
    done

    if [[ "$installed" == "" ]]; then
      @milpa.fail "No directory found for writing completion script (tried /etc/bash_completion.d and /usr/local/etc/bash_completion.d)"
    fi

    ;;
  *zsh)
    @milpa.log info "ZSH detected"

    # use zsh -c and source ~/.zshrc so we can read existing config without necessarily messing
    # with macos' session restoration and control character shenanigans
    # shellcheck disable=2016
    dst=$(zsh -c 'source ~/.zshrc 2>/dev/null; printf "\n%s" "${${fpath[@]:#$HOME/*}[1]}"' 2>/dev/null | tail -n 1) || @milpa.fail "Unable to locate an fpath to install completions to"
    @milpa.log info "Installing completions to $dst"
    if [[ -w "$dst" ]]; then
      [[ -f "$dst" ]] || mkdir -pv "$dst"
      "$MILPA_COMPA" __generate_completions zsh > "${dst}/_milpa" || @milpa.fail "Could not install completions"
    else
      @milpa.log warning "$dst does not look writeable for user $USER, using sudo"
      [[ -f "$dst" ]] || sudo mkdir -pv "$dst"
      "$MILPA_COMPA" __generate_completions zsh | sudo tee "${dst}/_milpa" >/dev/null || @milpa.fail "Could not install completions"
    fi

    if ! zsh -c "source ~/.zshrc 2>/dev/null; command -v compinit >/dev/null" >/dev/null 2>&1; then
      caveats='compinit has not been loaded into this shell, enable it by running:

echo "autoload -U compinit; compinit" >> ~/.zshrc

then reloading your shell'
    fi
    ;;
  *fish)
    @milpa.log info "Fish detected"
    dst="$HOME/.config/fish/completions/milpa.fish"
    "$MILPA_COMPA" __generate_completions fish > "$dst" || @milpa.fail "Could not install completions"
  ;;
  *)
    @milpa.fail "No completion script found for shell $SHELL"
esac

@milpa.log complete "Shell completion added for $SHELL successfully"
[[ "$caveats" == "" ]] || @milpa.log warning "$caveats"
exit 0
