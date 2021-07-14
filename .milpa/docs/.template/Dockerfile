FROM klakegg/hugo:ext-alpine

RUN apk add git --no-cache && \
  cd /src && \
  git config --global init.defaultBranch main && \
  git init . && \
  hugo new site . --force && \
  git submodule add https://github.com/google/docsy.git themes/docsy && \
  git submodule update --init --recursive && \
  mkdir -p /src/content

COPY package.json /src/
COPY config.toml /src/
COPY _variables_project.scss /src/assets/scss/
