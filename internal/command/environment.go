// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/runtime"
)

func (cmd *Command) EnvironmentMap() map[string]string {
	return map[string]string{
		_c.OutputCommandName: cmd.FullName(),
		_c.OutputCommandKind: string(cmd.Meta.Kind),
		_c.OutputCommandRepo: cmd.Meta.Repo,
		_c.OutputCommandPath: cmd.Meta.Path,
	}
}

func (cmd *Command) ToEval(args []string) string {
	output := []string{}
	for name, value := range runtime.EnvironmentMap() {
		output = append(output, fmt.Sprintf("export %s=%s", name, shellescape.Quote(value)))
	}

	for name, value := range cmd.EnvironmentMap() {
		output = append(output, fmt.Sprintf("export %s=%s", name, shellescape.Quote(value)))
	}

	cmd.Options.ToEnv(cmd, &output, "export ")
	cmd.Arguments.ToEnv(cmd, &output, "export ")

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	output = append(output, "set -- "+strings.Join(args, " "))

	return strings.Join(output, "\n")
}

func (cmd *Command) Env(seed []string) []string {
	for name, value := range runtime.EnvironmentMap() {
		seed = append(seed, fmt.Sprintf("%s=%s", name, shellescape.Quote(value)))
	}

	for name, value := range cmd.EnvironmentMap() {
		seed = append(seed, fmt.Sprintf("%s=%s", name, value))
	}

	cmd.Options.ToEnv(cmd, &seed, "")
	cmd.Arguments.ToEnv(cmd, &seed, "")

	return seed
}
