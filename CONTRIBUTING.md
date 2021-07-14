# Contributing to `milpa`

## Local setup

```sh
# setup git hooks and instal golang dependencies
# if using asdf-vm, install golang if needed
make setup
```

## Releasing

```sh
export MILPA_PATH="$(pwd)/internal"
milpa release create [major|minor|patch] [--pre [alpha|beta|rc]]
```
