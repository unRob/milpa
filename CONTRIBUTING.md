# Contributing to `milpa`

## Local setup

```sh
# `milpa` is self hosted-ish. You need `milpa` available to run milpa commands
curl -L https://milpa.dev/install.sh | bash -
git clone git@github.com:/unRob/milpa.git
cd milpa
MILPA_ROOT=$(pwd) ./milpa dev setup
```

## General flow

```sh
eval "$(milpa dev env)"
# then you're good to
milpa dev lint [go|shell]
milpa dev test [integration|unit]
# or all in one go
milpa dev ci
```

Remember to add notes to the changelog if making user-facing changes!

```sh
milpa cl add [breaking-change|bug|deprecation|feature|improvement|note] MESSAGE
```

## Releasing

```sh
milpa release create [--pre [alpha|beta|rc]]
```
