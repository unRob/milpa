SHELL := /usr/bin/env bash
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
	git config core.hooksPath $(shell git rev-parse --show-toplevel)/bin/hooks
	go get -u gotest.tools/gotestsum
	go mod tidy


test:
	gotestsum --format short -- ./...

lint:
	golangci-lint run
	shellcheck milpa .milpa/**/*.sh

compa: compa.go go.mod go.sum internal/*
	go build -ldflags "-s -w -X main.version=${MILPA_VERSION}" -o compa

