# Contributing to `milpa`

## Local setup

```sh
# setup git hooks and instal golang dependencies
# if using asdf-vm, install golang if needed
make setup
```

Remember to add notes to the changelog if making user-facing changes!

```sh
milpa cl add [breaking-change|bug|deprecation|feature|improvement|note] MESSAGE
```

## Releasing

```sh
milpa release create [--pre [alpha|beta|rc]]
```
