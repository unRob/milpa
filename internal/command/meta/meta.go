// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package meta

import (
	"path/filepath"
	"strings"

	"github.com/unrob/milpa/internal/command/kind"
	_c "github.com/unrob/milpa/internal/constants"
)

type Meta struct {
	// Path is the filesystem path to this command
	Path string `json:"path" yaml:"path"`
	// Repo is the filesystem path to this repo, including /.milpa
	Repo string `json:"repo" yaml:"repo"`
	// Name is a list of words naming this command
	Name []string `json:"name" yaml:"name"`
	// Kind can be executable (a binary or executable file), source (.sh file), or virtual (a sub-command group)
	Kind   kind.Kind `json:"kind" yaml:"kind"`
	Shell  string    `json:"shell" yaml:"shell"`
	Issues []error
}

func ForPath(path string, repo string) (meta Meta) {
	var name string
	if strings.HasSuffix(path, ".yaml") {
		name = filepath.Dir(path)
		name = strings.TrimPrefix(name, repo+"/")
		name = strings.TrimPrefix(name, _c.RepoCommandFolderName+"/")
		meta.Path = name
		meta.Kind = kind.Virtual
	} else {
		meta.Path = path
		extension := filepath.Ext(path)
		name = strings.TrimSuffix(path, extension)
		name = strings.TrimPrefix(name, repo+"/")
		name = strings.TrimPrefix(name, _c.RepoCommandFolderName+"/")

		switch extension {
		case ".zsh":
			meta.Kind = kind.Posix
			meta.Shell = "zsh"
		case ".sh", ".bash":
			meta.Kind = kind.Posix
			meta.Shell = "bash"
		default:
			meta.Kind = kind.Executable
		}
	}

	meta.Repo = repo
	meta.Name = strings.Split(name, "/")
	meta.Issues = []error{}

	return meta
}

func (meta *Meta) ParsingErrors() []error {
	return meta.Issues
}
