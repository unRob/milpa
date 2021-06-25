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
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CommandSetArgument struct {
	From struct {
		SubCommand string `yaml:"subcommand"`
	} `yaml:"from"`
	Values         []string
	computedValues *[]string
}

func (csa *CommandSetArgument) Resolve() ([]string, error) {
	values := []string{}
	if csa.From.SubCommand != "" {
		if csa.computedValues == nil {
			logrus.Debugf("executing sub command %s", csa.From.SubCommand)
			// milpa := fmt.Sprintf("%s/milpa", os.Getenv("MILPA_ROOT"))
			cmd := exec.Command("milpa", strings.Split(csa.From.SubCommand, " ")...)
			out, err := cmd.Output()
			if err != nil {
				logrus.Error(err)
				return values, err
			}

			val := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
			csa.computedValues = &val
		}
		values = *csa.computedValues
	} else if len(csa.Values) > 0 {
		return csa.Values, nil
	}

	return values, nil
}

type CommandArguments []CommandArgument

func (args *CommandArguments) ToEnv(dst *[]string, actual []string) error {
	for idx, arg := range *args {
		envName := fmt.Sprintf("MILPA_ARG_%s", strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_")))

		if idx >= len(actual) {
			if arg.Required {
				return fmt.Errorf("missing argument: %s", arg.Name)
			}
			logrus.Debugf("Skipping arg parsing for %s", arg.Name)
			value := ""
			if arg.Default != "" {
				value = arg.Default
			}
			*dst = append(*dst, fmt.Sprintf("export %s=%s", envName, value))
			continue
		}

		var value string
		if arg.Variadic {
			values := []string{}
			for _, va := range actual[idx:] {
				values = append(values, shellescape.Quote(va))
			}
			value = fmt.Sprintf("( %s )", strings.Join(values, " "))

		} else {
			value = shellescape.Quote(actual[idx])
		}

		if arg.Set != nil {
			values, err := arg.Set.Resolve()
			if err != nil {
				return err
			}
			found := false
			for _, validValue := range values {
				if value == validValue {
					found = true
					break
				}
			}

			if !found {
				return BadArguments{
					fmt.Sprintf("invalid value for %s: %s. Valid values are %s", arg.Name, value, strings.Join(values, ", ")),
				}
			}
		}
		logrus.Debugf("arg parsing for %s, %d, %s", arg.Name, idx, value)

		*dst = append(*dst, fmt.Sprintf("export %s=%s", envName, value))
	}

	return nil
}

func (args *CommandArguments) ToValidationFunction() func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	self := *args
	expectedArgLen := len(self)
	if expectedArgLen > 0 {
		return func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			argsCompleted := len(args)

			values := []string{}
			directive := cobra.ShellCompDirectiveDefault
			// logrus.Infof("argsCompleted: %d, expected: %d", argsCompleted, expectedArgLen)
			if argsCompleted < expectedArgLen {
				// el usuario pide completar un arg que aun esperamos
				arg := self[argsCompleted]
				if arg.Set != nil {
					values, _ = arg.Set.Resolve()
				} else {
					directive = cobra.ShellCompDirectiveError
				}
			}

			if toComplete != "" {
				filtered := []string{}
				for _, value := range values {
					if strings.HasPrefix(value, toComplete) {
						filtered = append(filtered, value)
					}
				}
				values = filtered
			}

			return values, directive
		}
	}
	return nil
}

type CommandArgument struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Default     string              `yaml:"default"`
	Set         *CommandSetArgument `yaml:"set"`
	Variadic    bool                `yaml:"variadic"`
	Required    bool                `yaml:"required"`
}

func (cmdarg *CommandArgument) Validates() bool {
	return cmdarg.Set != nil
}

func (cmdarg *CommandArgument) ToDesc() string {
	spec := strings.ToUpper(cmdarg.Name)

	if !cmdarg.Required {
		spec = fmt.Sprintf("[%s]", spec)
	}

	if cmdarg.Variadic {
		spec = fmt.Sprintf("%s...", spec)
	}
	return spec
}

type ValueType string

const (
	ValueTypeDefault ValueType = ""
	ValueTypeString  ValueType = "string"
	ValueTypeBoolean ValueType = "boolean"
)

type CommandOption struct {
	ShortName   string      `yaml:"short-name"`
	Type        ValueType   `yaml:"type"`
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
}
