// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/exec"
	"github.com/spf13/cobra"
)

func MilpaComplete(cmd *command.Command, currentValue string, config string) (values []string, flag cobra.ShellCompDirective, err error) {
	cmdLine, err := cmd.ResolveTemplate(config, currentValue)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError, err
	}

	args := append([]string{"milpa"}, strings.Split(cmdLine, " ")...)
	envMap := EnvironmentMap(cmd)
	env := os.Environ()
	for k, v := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	values, flag, err = exec.Exec(cmd.FullName(), args, env, 5*time.Second)
	if err != nil {
		return nil, flag, err
	}

	return
}

func init() {
	command.RegisterValueSource("milpa", MilpaComplete)
}
