// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package kind

import (
	"fmt"
	"os/exec"
)

type Kind string

const (
	Unknown     Kind = ""
	Executable  Kind = "executable"
	Source      Kind = "source"
	Virtual     Kind = "virtual"
	Root        Kind = "root"
	ShellScript Kind = "shell-script"
)

type Shell string

const (
	ShellUnknown Shell = ""
	ShellZSH     Shell = "zsh"
	ShellBash    Shell = "bash"
	ShellFish    Shell = "fish"
)

func (s Shell) Path() (string, error) {
	shell, err := exec.LookPath(string(s))
	if err != nil {
		return shell, fmt.Errorf("could not find an executable for %s: %s", s, err)
	}

	return shell, nil
}
