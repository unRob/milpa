FROM klakegg/hugo:ext-alpine

RUN apk add git --no-cache && \
  cd /src && \
  hugo new site . --force && \
  mkdir -p /src/content/themes

COPY theme /src/themes/cli
COPY config.toml /src/
