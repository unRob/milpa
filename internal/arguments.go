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
	"github.com/spf13/pflag"
)

type Arguments []Argument

func (args *Arguments) ToEnv(dst *[]string, actual []string) error {
	for idx, arg := range *args {
		envName := fmt.Sprintf("MILPA_ARG_%s", strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_")))

		if idx >= len(actual) {
			if arg.Required {
				return fmt.Errorf("missing argument: %s", strings.ToUpper(arg.Name))
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

		if arg.Validates() {
			values, err := arg.Resolve()
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

func (args *Arguments) Validate(cc *cobra.Command, supplied []string) error {
	for idx, arg := range *args {
		argumentProvided := idx < len(supplied)
		if arg.Required && !argumentProvided {
			return BadArguments{fmt.Sprintf("Missing argument for %s", strings.ToUpper(arg.Name))}
		}

		if !argumentProvided {
			continue
		}
		current := supplied[idx]

		if arg.Validates() {
			values, err := arg.Resolve()
			if err != nil {
				return err
			}
			found := false
			for _, value := range values {
				if value == current {
					found = true
					break
				}
			}

			if !found {
				return BadArguments{fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(values, ", "))}
			}
		}
	}

	return nil
}

func (args *Arguments) CompletionFunction() func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
				if arg.Validates() {
					values, _ = arg.Resolve()
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

type Argument struct {
	Name             string   `yaml:"name" validate:"required,excludesall=!$\\/%^@#?:'\""`
	Description      string   `yaml:"description" validate:"required"`
	Default          string   `yaml:"default" validate:"excluded_with=Required"`
	Variadic         bool     `yaml:"variadic"`
	Required         bool     `yaml:"required" validate:"excluded_with=Default"`
	ValuesSubCommand string   `yaml:"values-subcommand" validate:"excluded_with=Values"`
	Values           []string `yaml:"values" validate:"excluded_with=ValuesSubCommand"`
	computedValues   *[]string
}

func (arg *Argument) Validates() bool {
	return len(arg.Values) > 0 || arg.ValuesSubCommand != ""
}

func (arg *Argument) ToDesc() string {
	spec := strings.ToUpper(arg.Name)

	if !arg.Required {
		spec = fmt.Sprintf("[%s]", spec)
	}

	if arg.Variadic {
		spec = fmt.Sprintf("%s...", spec)
	}
	return spec
}

func (arg *Argument) Resolve() ([]string, error) {
	values := []string{}
	if arg.ValuesSubCommand != "" {
		if arg.computedValues == nil {
			logrus.Debugf("executing sub command %s", arg.ValuesSubCommand)
			// milpa := fmt.Sprintf("%s/milpa", os.Getenv("MILPA_ROOT"))
			cmd := exec.Command("milpa", strings.Split(arg.ValuesSubCommand, " ")...) // #nosec G204
			out, err := cmd.Output()
			if err != nil {
				logrus.Error(err)
				return values, err
			}

			val := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
			arg.computedValues = &val
		}
		values = *arg.computedValues
	} else if len(arg.Values) > 0 {
		values = arg.Values
	}

	return values, nil
}

type ValueType string

const (
	ValueTypeDefault ValueType = ""
	ValueTypeString  ValueType = "string"
	ValueTypeBoolean ValueType = "bool"
)

type Option struct {
	ShortName        string      `yaml:"short-name"`
	Type             ValueType   `yaml:"type" validate:"omitempty,oneof=string bool"`
	Description      string      `yaml:"description" validate:"required"`
	Default          interface{} `yaml:"default"`
	ValuesSubCommand string      `yaml:"values-subcommand" validate:"excluded_with=Values"`
	Values           []string    `yaml:"values" validate:"excluded_with=ValuesSubCommand"`
	computedValues   *[]string
}

func (opt *Option) Validates() bool {
	return len(opt.Values) > 0 || opt.ValuesSubCommand != ""
}

func (opt *Option) Resolve() ([]string, error) {
	values := []string{}
	if opt.ValuesSubCommand != "" {
		if opt.computedValues == nil {
			logrus.Debugf("executing sub command %s", opt.ValuesSubCommand)
			cmd := exec.Command("milpa", strings.Split(opt.ValuesSubCommand, " ")...) // #nosec G204
			out, err := cmd.Output()
			if err != nil {
				logrus.Error(err)
				return values, err
			}

			val := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
			opt.computedValues = &val
		}
		values = *opt.computedValues
	} else if len(opt.Values) > 0 {
		values = opt.Values
	}

	return values, nil
}

func (opt *Option) ValidationFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	values, err := opt.Resolve()
	if err != nil {
		return values, cobra.ShellCompDirectiveError
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
	return values, cobra.ShellCompDirectiveDefault
}

type Options map[string]*Option

func (opts *Options) ToEnv(dst *[]string, flags *pflag.FlagSet) (err error) {
	errors := []string{}
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
			oopts := *opts
			opt, ok := oopts[name]
			if ok && opt.Validates() {
				values, verr := opt.Resolve()
				if verr != nil {
					errors = append(errors, err.Error())
					return
				}
				found := false
				for _, validValue := range values {
					if value == validValue {
						found = true
						break
					}
				}

				if !found {
					errors = append(errors,
						fmt.Sprintf("invalid value for --%s: %s. Valid values are %s", name, value, strings.Join(values, ", ")),
					)
				}
			}
			logrus.Debugf("flag %s is a %s", f.Name, f.Value.Type())
		}

		value = shellescape.Quote(value)
		*dst = append(*dst, fmt.Sprintf("export %s=%s", envName, value))
	})

	if len(errors) > 0 {
		return BadArguments{strings.Join(errors, ". ")}
	}
	return nil
}
