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
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
)

func contains(haystack []string, needle string) bool {
	for _, validValue := range haystack {
		if needle == validValue {
			return true
		}
	}
	return false
}

// Arguments is an ordered list of Argument.
type Arguments []*Argument

func (args *Arguments) AllKnown() map[string]string {
	col := map[string]string{}
	for _, arg := range *args {
		col[arg.Name] = arg.ToString(false)
	}
	return col
}

// ToEnv writes shell variables to dst.
func (args *Arguments) ToEnv(cmd *Command, dst *[]string) {
	for _, arg := range *args {
		envName := fmt.Sprintf("%s%s", _c.OutputPrefixArg, arg.EnvName())

		*dst = append(*dst, fmt.Sprintf("export %s=%s", envName, arg.ToString(true)))
	}
}

func (args *Arguments) Parse(supplied []string) {
	for idx, arg := range *args {
		argumentProvided := idx < len(supplied)

		if !argumentProvided {
			if arg.Default != "" {
				arg.provided = &[]string{arg.Default}
			}
			continue
		}

		if arg.Variadic {
			values := []string{}
			for _, va := range supplied[idx:] {
				values = append(values, shellescape.Quote(va))
			}
			arg.SetValue(values)
		} else {
			arg.SetValue([]string{supplied[idx]})
		}
	}
}

func (args *Arguments) AreValid(cmd *Command) error {
	for _, arg := range *args {
		if err := arg.Validate(cmd); err != nil {
			return err
		}
	}

	return nil
}

// CompletionFunction is called by cobra when asked to complete arguments.
func (args *Arguments) CompletionFunction(command *Command) func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	self := *args
	expectedArgLen := len(self)
	hasVariadicArg := expectedArgLen > 0 && self[len(self)-1].Variadic
	if expectedArgLen > 0 {
		return func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			argsCompleted := len(args)
			command.Options.Parse(cc.Flags())
			command.Arguments.Parse(cc.Flags())

			values := []string{}
			directive := cobra.ShellCompDirectiveDefault
			if argsCompleted < expectedArgLen || hasVariadicArg {
				var arg *Argument
				if hasVariadicArg && argsCompleted >= expectedArgLen {
					// completing a variadic argument
					arg = self[len(self)-1]
				} else {
					// completing regular argument (maybe variadic!)
					arg = self[argsCompleted]
				}

				if arg.Values != nil {
					values, directive, _ = arg.Resolve(command)
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

// Argument represents a single command-line argument.
type Argument struct {
	// Name is how this variable will be exposed to the underlying command.
	Name string `yaml:"name" validate:"required,excludesall=!$\\/%^@#?:'\""`
	// Description is what this argument is for.
	Description string `yaml:"description" validate:"required"`
	// Default is the default value for this argument if none is provided.
	Default string `yaml:"default" validate:"excluded_with=Required"`
	// Variadic makes an argument a list of all values from this one on.
	Variadic bool `yaml:"variadic"`
	// Required raises an error if an argument is not provided.
	Required bool `yaml:"required" validate:"excluded_with=Default"`
	// Values (DEPRECATED, renaming to Values) describes autocompletion and validation for an argument
	Values   *ValueSource `yaml:"values" validate:"omitempty"`
	provided *[]string
}

func (arg *Argument) EnvName() string {
	return strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_"))
}

func (arg *Argument) SetValue(value []string) {
	arg.provided = &value
}

func (arg *Argument) IsKnown() bool {
	return arg.provided != nil
}

func (arg *Argument) ToString(asShell bool) string {
	value := arg.Default
	if arg.IsKnown() {
		if arg.Variadic {
			if asShell {
				values := []string{}
				for _, va := range *arg.provided {
					values = append(values, shellescape.Quote(va))
				}
				value = fmt.Sprintf("( %s )", strings.Join(values, " "))
			} else {
				value = strings.Join(*arg.provided, " ")
			}
		} else {
			vals := *arg.provided
			value = vals[0]
		}
	}

	if !arg.Variadic && asShell {
		value = shellescape.Quote(value)
	}

	return value
}

func (arg *Argument) Validate(cmd *Command) error {

	if !arg.IsKnown() {
		if arg.Required {
			return BadArguments{fmt.Sprintf("Missing argument for %s", strings.ToUpper(arg.Name))}
		}

		return nil
	}

	if !arg.Validates() {
		return nil
	}

	validValues, _, err := arg.Resolve(cmd)
	if err != nil {
		return err
	}

	if arg.Variadic {
		for _, current := range *arg.provided {
			if !contains(validValues, current) {
				return BadArguments{fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(validValues, ", "))}
			}
		}
	} else {
		current := arg.ToString(false)
		if !contains(validValues, current) {
			return BadArguments{fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(validValues, ", "))}
		}
	}

	return nil
}

// Validates tells if the user-supplied value needs validation.
func (arg *Argument) Validates() bool {
	return arg.Values != nil && arg.Values.Validates()
}

// ToDesc prints out the description of an argument for help and docs.
func (arg *Argument) ToDesc() string {
	spec := arg.EnvName()
	if arg.Variadic {
		spec = fmt.Sprintf("%s...", spec)
	}

	if !arg.Required {
		spec = fmt.Sprintf("[%s]", spec)
	}

	return spec
}

// Resolve returns autocomplete values for an argument.
func (arg *Argument) Resolve(command *Command) (values []string, flag cobra.ShellCompDirective, err error) {
	if arg.Values != nil {
		values, flag, err = arg.Values.Resolve(command)
		if err != nil {
			flag = cobra.ShellCompDirectiveError
			return
		}
	}

	return
}
