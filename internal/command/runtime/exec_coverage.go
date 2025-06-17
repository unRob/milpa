// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//go:build coverage

package runtime

import (
	"os"
	"os/exec"
)

func fork(command string, args, cmdEnv []string) error {
	cmd := &exec.Cmd{
		Path:   command,
		Args:   args,
		Env:    cmdEnv,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
