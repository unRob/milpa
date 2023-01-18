// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"strings"

	_c "github.com/unrob/milpa/internal/constants"
)

type Kind string

const (
	KindUnknown    Kind = ""
	KindExecutable Kind = "executable"
	KindSource     Kind = "source"
	KindVirtual    Kind = "virtual"
	KindRoot       Kind = "root"
)

type Meta struct {
	// Path is the filesystem path to this command
	Path string `json:"path" yaml:"path"`
	// Repo is the filesystem path to this repo, including /.milpa
	Repo string `json:"repo" yaml:"repo"`
	// Name is a list of words naming this command
	Name []string `json:"name" yaml:"name"`
	// Kind can be executable (a binary or executable file), source (.sh file), or virtual (a sub-command group)
	Kind   Kind `json:"kind" yaml:"kind"`
	issues []error
}

func metaForPath(path string, repo string) (meta Meta) {
	meta.Path = path
	meta.Repo = repo
	name := strings.TrimSuffix(path, ".sh")
	name = strings.TrimPrefix(name, repo+"/"+_c.RepoCommandFolderName+"/")
	meta.Name = strings.Split(name, "/")
	meta.issues = []error{}

	if strings.HasSuffix(path, ".sh") {
		meta.Kind = "source"
	} else {
		meta.Kind = "exec"
	}

	return
}
