// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//go:build !coverage

package runtime

import "syscall"

// this is used on release builds because I truly want to replace the current process with whatever's next
// however, fork and exec seeminlgy breaks coverage reporting.
func fork(command string, args, cmdEnv []string) error {
	return syscall.Exec(command, args, cmdEnv)
}
