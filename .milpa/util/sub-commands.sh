#!/usr/bin/env bash

_MILPA_COMMAND_CACHE="$HOME/.config/milpa/commands.json"

#shellcheck disable=SC2120
function __find_commands() {
  local query_prefix args subpath;
  if [[ -n "${1+x}" ]]; then
    query_prefix="/${1}"
  fi
  subpath="/.milpa/commands$query_prefix"
  IFS=':' read -r -a args <<< "${MILPA_PATH//:/$subpath:}$subpath"
  args+=(
    -type f
    -name '*.sh' -o ! -name '*.*'
  )

  if [[ -n "${1+x}" ]]; then
    args+=( -maxdepth 1 )
  fi

  find "${args[@]}" 2>/dev/null | sort
}

function __build_command_list () {
  local path
  mkdir -p "$(dirname "$_MILPA_COMMAND_CACHE")"
  __find_commands | while read -r path; do
    [[ -d "$path" ]] && continue

    fname="$(basename "$path")"
    [[ "$fname" == _env* ]] && continue

    dir=$(dirname "$path")
    name=${fname%%.*}
    full_name="${dir##*.milpa\/commands}/$name"
    docs="$dir/${name}.md"
    description=$(awk '/^([^#])/ { print $0; exit; }' "$docs" 2>/dev/null)

    echo "${full_name:1}!-milpa-!$path!-milpa-!$description"
  done | jq -Rs '
  split("\n") |
  sort |
  map(select(length > 0) | split("!-milpa-!") | {
    key: (.[0] | split("/")),
    value: {
      fullName: (.[0] | gsub("/"; " ")),
      description: .[2],
      path: .[1],
      package: (.[1] | split("/.milpa")[0]),
      _command: true,
    }
  }) |
  # asdf
  reduce .[] as $cmd ({};
    . | setpath($cmd.key; $cmd.value)
  )' > "$_MILPA_COMMAND_CACHE"
}

function _sub_commands_descriptions () {
  jq -r --arg "root" "${1}" '
    if $root != "" then getpath([$root]) else . end |
    [
      .. |
      select(._command?) |
      "\(.fullName) - \(.description)"
    ][]
  ' "$_MILPA_COMMAND_CACHE"
}

function _sub_commands_packages () {
  jq -r'
    [
      .. |
      select(._command?) |
      "\(.fullName) - \(.description)"
    ][]
  ' "$_MILPA_COMMAND_CACHE"
}

function _sub_commands_find () {
  jq -r \
    --slurpfile commands "$_MILPA_COMMAND_CACHE" \
    --null-input \
    '
  $ARGS.positional as $args |
  $args | length as $total |
  def find_command($idx):
    $args[0:$total - $idx] as $query |
    # "looking at depth \($idx) \($query)" | debug |
    $commands[0] | getpath($query) as $cmd |
    if $cmd.path // null != null then
      # "found \($cmd.path)" | debug |
      [($query | join(" ")), $cmd.path, $cmd.package, $total-$idx]
    else
      # "nothing for \($query)" | debug |
      if $idx + 1 <= $total then
        find_command($idx + 1)
      else
        "Unknown command\n" | halt_error(2)
      end
    end;
  find_command(0)[]
  ' \
  --args -- "${@}"
}

#[[ ! -f "$_MILPA_COMMAND_CACHE" ]] &&
 __build_command_list
