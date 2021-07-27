SHELL := /usr/bin/env bash -O globstar
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
.PHONY: compa setup test lint clean
RELEASE_TARGET ?= dist
TARGET_MACHINES = linux-amd64 linux-arm64 linux-arm linux-mips darwin-amd64 darwin-arm64
TARGET_ARCHIVES = $(addsuffix .tgz,$(addprefix $(RELEASE_TARGET)/packages/milpa-,$(TARGET_MACHINES)))

ifneq ($(ASDF_DIR), "")
setup-golang:
	$(info Installing golang version with asdf)
	asdf plugin list | grep golang >/dev/null || asdf plugin add golang
	asdf list golang | grep -f <(cut -d" " -f 2 .tool-versions) >/dev/null || asdf install
	$(shell asdf reshim)
else
setup-golang:
	$(warning Please make sure golang version $(shell cut -d" " -f 2 .tool-versions) is installed)
endif

setup: setup-golang
	$(info Configuring git hooks)
	git config core.hooksPath $(shell git rev-parse --show-toplevel)/internal/bin/hooks
	$(info Installing dev go packages)
	go get -u gotest.tools/gotestsum
	go get -u github.com/hashicorp/go-getter/cmd/go-getter
	go mod tidy

# every day usage
test:
	gotestsum --format short -- ./...

lint:
	golangci-lint run
	shellcheck milpa bootstrap.sh .milpa/**/*.sh repos/internal/**/*.sh

compa: compa.go go.mod go.sum internal/*
	go build -ldflags "-s -w -X main.version=${MILPA_VERSION}" -o compa

# Releasing
$(RELEASE_TARGET)/packages/milpa-%.tgz: compa.go go.mod go.sum internal/*.go
	mkdir -p $(basename $(subst packages,tmp,$@))/milpa

	GOOS=$(firstword $(subst -, ,$*)) GOARCH=$(lastword $(subst -, ,$*)) \
		go build -ldflags "-s -w -X main.version=${MILPA_VERSION}" -trimpath \
		-o $(basename $(subst packages,tmp,$@))/milpa/compa
	upx --no-progress -9 $(basename $(subst packages,tmp,$@))/milpa/compa

	cp -r ./milpa ./.milpa LICENSE.txt README.md CHANGELOG.md $(basename $(subst packages,tmp,$@))/milpa
	mkdir -p $(dir $@)
	tar -czf $@ -C $(basename $(subst packages,tmp,$@)) milpa
	openssl dgst -sha256 $@ | awk '{print $$2}' > $(subst .tgz,.shasum,$@)

$(RELEASE_TARGET)/packages: $(TARGET_ARCHIVES)
	rm -rf $(RELEASE_TARGET)/tmp
	$(info Built for $(TARGET_ARCHIVES))

clean:
	rm -rf $(RELEASE_TARGET)
