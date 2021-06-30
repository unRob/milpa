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
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var customNames map[string]string = map[string]string{"no-color": "NO_COLOR", "silent": "MILPA_SILENT", "verbose": "MILPA_VERBOSE"}

func setEnvForOpts(env *[]string, flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if name == "help" {
			return
		}
		envName := ""
		value := f.Value.String()

		if cname, ok := customNames[name]; ok {
			if value == "false" {
				return
			}
			envName = cname
		} else {
			envName = fmt.Sprintf("MILPA_OPT_%s", strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
		}

		switch f.Value.Type() {
		case "bool":
			if val, err := flags.GetBool(f.Name); err == nil && !val {
				value = ""
			} else {
				value = "true"
			}
		default:
			logrus.Debugf("flag %s is a %s", f.Name, f.Value.Type())
		}

		value = shellescape.Quote(value)
		*env = append(*env, fmt.Sprintf("export %s=%s", envName, value))
	})
}

func (cmd *Command) ToEval(args []string, flags *pflag.FlagSet) (string, error) {
	envVars := []string{
		fmt.Sprintf("export MILPA_COMMAND_NAME=%s", shellescape.Quote(cmd.FullName())),
		fmt.Sprintf("export MILPA_COMMAND_KIND=%s", shellescape.Quote(cmd.Meta.Kind)),
		fmt.Sprintf("export MILPA_COMMAND_REPO=%s", shellescape.Quote(cmd.Meta.Repo)),
		fmt.Sprintf("export MILPA_COMMAND_PATH=%s", shellescape.Quote(cmd.Meta.Path)),
	}

	setEnvForOpts(&envVars, flags)

	logrus.Debugf("Printing environment for args: %v", args)

	err := cmd.Arguments.ToEnv(&envVars, args)
	if err != nil {
		return "", err
	}

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	envVars = append(envVars, "set -- "+strings.Join(args, " "))

	return strings.Join(envVars, "\n"), nil
}
