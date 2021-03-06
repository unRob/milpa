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
		_c.EnvVarMilpaPath:       strings.Join(runtime.MilpaPath, ":"),
		_c.EnvVarMilpaPathParsed: "true",
		_c.OutputCommandName:     cmd.FullName(),
		_c.OutputCommandKind:     string(cmd.Meta.Kind),
		_c.OutputCommandRepo:     cmd.Meta.Repo,
		_c.OutputCommandPath:     cmd.Meta.Path,
	}
}

func (cmd *Command) ToEval(args []string) string {
	output := []string{}
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
	for name, value := range cmd.EnvironmentMap() {
		seed = append(seed, fmt.Sprintf("%s=%s", name, shellescape.Quote(value)))
	}

	cmd.Options.ToEnv(cmd, &seed, "")
	cmd.Arguments.ToEnv(cmd, &seed, "")

	return seed
}
