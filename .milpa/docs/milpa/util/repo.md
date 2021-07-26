---
description: Functions related to milpa repositories
---

The `repo` util contains shell functions related to milpa repositories

## Functions

### `@milpa.repo.current_path`

`@milpa.repo.current_path`

Returns the nearest `.milpa` repo from the current working directory where `milpa` was invoked from. It returns `2` if it reaches the `/` directory without finding a `.milpa` folder in its path.
