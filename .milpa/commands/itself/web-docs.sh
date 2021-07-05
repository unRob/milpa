#!/usr/bin/env bash

mkdir -p "$MILPA_ROOT/milpa.dev/content/docs"
for path in ${MILPA_PATH_ARR[*]}; do
  src="${path}/.milpa/docs"
  [[ ! -d "$src" ]] && continue

  @milpa.log info "copying docs from $src"
  cd "$src"
  cp -vr ./* "$MILPA_ROOT/milpa.dev/content/docs/"
  cd -
done

find "$MILPA_ROOT/milpa.dev/content/docs" -name "index.md" | while read -r index; do
  mv "$index" "${index//index.md/_index.md}"
done

find "$MILPA_ROOT/milpa.dev/content/docs" -name "*.md" | while read -r doc; do
  sed -i '' -E 's|\(/.milpa\/|(/|g; s|/index.md|/|g; s|\.md([\)#])|\1|g; s|!milpa!|'"$MILPA_NAME"'|g' "$doc"
done

@milpa.log info "generating command docs"
MILPA_PLAIN_HELP=enabled "$MILPA_COMPA" __generate_documentation "$MILPA_ROOT/milpa.dev/_tmp" || @milpa.fail "Could not generate command documentation"
mv "$MILPA_ROOT/milpa.dev/_tmp/milpa" "$MILPA_ROOT/milpa.dev/content/commands" || @milpa.fail "Failed to move commands"

@milpa.log info "Launching hugo website generator"
docker run --rm -it -p 1313:1313 -v "$(pwd)/milpa.dev:/src" milpa-docs server
