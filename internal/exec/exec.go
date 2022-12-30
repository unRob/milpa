// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package exec

import (
	"bytes"
	"context"
	"fmt"
	os_exec "os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal/errors"
)

// ExecFunc is replaced in tests.
var ExecFunc = WithSubshell

func WithSubshell(ctx context.Context, env []string, executable string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := os_exec.CommandContext(ctx, executable, args...) // #nosec G204
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Env = env
	return stdout, stderr, cmd.Run()
}

// Exec runs a subprocess and returns a list of lines from stdout.
func Exec(name string, args []string, env []string, timeout time.Duration) ([]string, cobra.ShellCompDirective, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	logrus.Debugf("executing %s", args)
	executable := args[0]
	args = args[1:]

	stdout, _, err := ExecFunc(ctx, env, executable, args...)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Sub-command timed out")
		logrus.Debugf("timeout running %s %s: %s", executable, args, stdout.String())
		return []string{}, cobra.ShellCompDirectiveError, fmt.Errorf("timed out resolving %s %s", executable, args)
	}

	if err != nil {
		logrus.Debugf("error running %s %s: %s", executable, args, err)
		return []string{}, cobra.ShellCompDirectiveError, errors.BadArguments{Msg: fmt.Sprintf("could not validate argument for command %s, ran <%s %s> failed: %s", name, executable, strings.Join(args, " "), err)}
	}

	logrus.Debugf("done running %s %s: %s", executable, args, stdout.String())
	return strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n"), cobra.ShellCompDirectiveDefault, nil
}
