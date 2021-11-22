# Contributing to `milpa`

## Local setup

```sh
# `milpa` is self hosted-ish.
# You'll need `milpa` available to run milpa commands
curl -L https://milpa.dev/install.sh | bash -
git clone git@github.com:/unRob/milpa.git
cd milpa
# temporarily use stable compa to run HEAD milpa
ln -sfv "$(dirname "$(which milpa)")/compa" $(pwd)/compa
# do all sorts of crazy setup, including bulding compa
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
# build local compa binaries with
milpa dev build
```

Remember to add notes to the changelog if making user-facing changes!

```sh
milpa cl add [breaking-change|bug|deprecation|feature|improvement|note] MESSAGE
```

## Releasing

```sh
milpa release create [--pre [alpha|beta|rc]]
```
