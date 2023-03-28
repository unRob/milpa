---
description: Functions related to user input
---

The `user-input` util contains shell functions related to the creation and cleanup of temporary files and directories.

## Functions

### `@milpa.ask`

`@milpa.ask PROMPT [DEFAULT]`

Prompts the user to enter a value. It will show the prompt, followed by the default value with no additional formatting. It outputs the entered value on success and asks again if no value was entered (and no default is set).

```sh
#!/usr/bin/env bash
@milpa.load_util user-input

# prompt the user to enter their favorite food
food=$(@milpa.ask "What's your favorite ingredient for food?" "tortillas")
# shows:
# What's your favorite ingredient for food? [default: tortillas]

@milpa.log info "$food are the best!"
```

### `@milpa.confirm`

`@milpa.confirm [PROMPT]`

Prompts the user to press the `y` key to continue the execution of a script. Optionally, a `PROMPT` may be presented to the user.


```sh
#!/usr/bin/env bash
@milpa.load_util user-input

# prompt the user for confirmation
if @milpa.confirm "Do you wanna see something cool?"; then
  @milpa.log info "🧊"
else
  @milpa.log info "Maybe next time!"
fi
```

### `@milpa.select`

`@milpa.select "OPTION\n[...]"`

Prompts the user to select a number from a list of options printed to the screen.


```sh
#!/usr/bin/env bash
@milpa.load_util user-input

@milpa.log info "What is your favorite vitamin-T meal?"
meal="$(cat <<<OPTIONS | @milpa.select
tacos
tamales
tlayudas
tortas
tostadas
OPTIONS
)"

@milpa.log success "Definitely, $meal are delicious"
```

