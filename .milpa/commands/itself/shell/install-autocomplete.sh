#!/usr/bin/env bash

case "$SHELL" in
  *bash)
    if [[ -d /etc/bash_completion.d ]]; then
      set -x
      "$MILPA_HELPER" __generate_completions bash > /etc/bash_completion.d/milpa
      set +x
    elif [[ -d /usr/local/etc/bash_completion.d ]]; then
      set -x
      "$MILPA_HELPER" __generate_completions bash > /usr/local/etc/bash_completion.d/milpa
      set +x
    else
      _fail "No directory found for writing completion script (tried /etc/bash_completion.d and /usr/local/etc/bash_completion.d)"
    fi
    ;;
  *zsh)
    $SHELL -i -c "command -v compinit >/dev/null" >/dev/null || _log warning <<EOF
compinit has not been loaded into this shell, enable it by running

echo "autoload -U compinit; compinit" >> ~/.zshrc
and reloading your shell
EOF
    # shellcheck disable=2016
    zsh -i -c '
dst="${${fpath[@]:#$HOME/*}[1]}"
set -ex
"$MILPA_HELPER" __generate_completions zsh > "${dst}/_milpa"
set +x' && _log warning "Please restart your shell"
    ;;
  *fish)
    set -ex
    "$MILPA_HELPER" __generate_completions fish > ~/.config/fish/completions/milpa.fish
    set +x
  ;;
  *)
    _fail "No completion script found for shell $SHELL"
esac

_log complete "Shell completion added for $SHELL successfully"
