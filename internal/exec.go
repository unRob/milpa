// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	os_exec "os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExecFunc is replaced in tests.
var ExecFunc = subshell

func subshell(ctx context.Context, env []string, executable string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := os_exec.CommandContext(ctx, executable, args...) // #nosec G204
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	return stdout, stderr, cmd.Run()
}

// Exec runs a subprocess and returns a list of lines from stdout.
func Exec(name string, args []string, timeout time.Duration) ([]string, cobra.ShellCompDirective, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	logrus.Debugf("executing %s", args)
	executable := args[0]
	args = args[1:]

	stdout, _, err := ExecFunc(ctx, os.Environ(), executable, args...)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Sub-command timed out")
		logrus.Debugf("timeout running %s %s: %s", executable, args, stdout.String())
		return []string{}, cobra.ShellCompDirectiveError, fmt.Errorf("timed out resolving %s %s", executable, args)
	}

	if err != nil {
		logrus.Debugf("error running %s %s: %s", executable, args, err)
		return []string{}, cobra.ShellCompDirectiveError, BadArguments{fmt.Sprintf("could not validate argument %s, sub-command <%s> failed: %s", name, args, err)}
	}

	logrus.Debugf("done running %s %s: %s", executable, args, stdout.String())
	return strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n"), cobra.ShellCompDirectiveDefault, nil
}
