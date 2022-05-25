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
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
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
func (args *Arguments) ToEnv(cmd *Command, dst *[]string, prefix string) {
	for _, arg := range *args {
		envName := fmt.Sprintf("%s%s", _c.OutputPrefixArg, arg.EnvName())

		*dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, envName, arg.ToString(true)))
	}
}

func (args *Arguments) Parse(supplied []string) {
	for idx, arg := range *args {
		argumentProvided := idx < len(supplied)

		if !argumentProvided {
			if arg.Default != nil {
				if arg.Variadic {
					defaultSlice := []string{}
					for _, valI := range arg.Default.([]interface{}) {
						defaultSlice = append(defaultSlice, valI.(string))
					}
					arg.provided = &defaultSlice
				} else {
					defaultString := arg.Default.(string)
					if defaultString != "" {
						arg.provided = &[]string{defaultString}
					}
				}
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

func (args *Arguments) AreValid() error {
	for _, arg := range *args {
		if err := arg.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// CompletionFunction is called by cobra when asked to complete arguments.
func (args *Arguments) CompletionFunction(cc *cobra.Command, provided []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	expectedArgLen := len(*args)
	values := []string{}
	directive := cobra.ShellCompDirectiveError

	if expectedArgLen > 0 {
		argsCompleted := len(provided)
		lastArg := (*args)[len(*args)-1]
		hasVariadicArg := expectedArgLen > 0 && lastArg.Variadic
		lastArg.Command.Options.Parse(cc.Flags())
		args.Parse(provided)

		directive = cobra.ShellCompDirectiveDefault
		if argsCompleted < expectedArgLen || hasVariadicArg {
			var arg *Argument
			if hasVariadicArg && argsCompleted >= expectedArgLen {
				// completing a variadic argument
				arg = lastArg
			} else {
				// completing regular argument (maybe variadic!)
				arg = (*args)[argsCompleted]
			}

			if arg.Values != nil {
				var err error
				arg.Values.Command = lastArg.Command
				arg.Command = lastArg.Command
				values, directive, err = arg.Resolve()
				if err != nil {
					return []string{err.Error()}, cobra.ShellCompDirectiveDefault
				}
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
	}

	return values, directive
}

// Argument represents a single command-line argument.
type Argument struct {
	// Name is how this variable will be exposed to the underlying command.
	Name string `yaml:"name" validate:"required,excludesall=!$\\/%^@#?:'\""`
	// Description is what this argument is for.
	Description string `yaml:"description" validate:"required"`
	// Default is the default value for this argument if none is provided.
	Default interface{} `yaml:"default" validate:"excluded_with=Required"`
	// Variadic makes an argument a list of all values from this one on.
	Variadic bool `yaml:"variadic"`
	// Required raises an error if an argument is not provided.
	Required bool `yaml:"required" validate:"excluded_with=Default"`
	// Values describes autocompletion and validation for an argument
	Values   *ValueSource `yaml:"values" validate:"omitempty"`
	Command  *Command     `json:"-" validate:"-"`
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
	var value string
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
	} else {
		if arg.Default != nil {
			if arg.Variadic {
				defaultSlice := []string{}
				for _, valI := range arg.Default.([]interface{}) {
					valStr := valI.(string)
					if asShell {
						valStr = shellescape.Quote(valStr)
					}
					defaultSlice = append(defaultSlice, valStr)
				}

				value = strings.Join(defaultSlice, " ")
				if asShell {
					value = fmt.Sprintf("( %s )", value)
				}
			} else {
				value = arg.Default.(string)
			}
		} else {
			value = ""
		}
	}

	if !arg.Variadic && asShell {
		value = shellescape.Quote(value)
	}

	return value
}

func (arg *Argument) Validate() error {

	if !arg.IsKnown() {
		if arg.Required {
			return errors.BadArguments{Msg: fmt.Sprintf("Missing argument for %s", strings.ToUpper(arg.Name))}
		}

		return nil
	}

	if !arg.Validates() {
		return nil
	}

	validValues, _, err := arg.Resolve()
	if err != nil {
		return err
	}

	if arg.Variadic {
		for _, current := range *arg.provided {
			if !contains(validValues, current) {
				return errors.BadArguments{Msg: fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(validValues, ", "))}
			}
		}
	} else {
		current := arg.ToString(false)
		if !contains(validValues, current) {
			return errors.BadArguments{Msg: fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(validValues, ", "))}
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
func (arg *Argument) Resolve() (values []string, flag cobra.ShellCompDirective, err error) {
	if arg.Values != nil {
		values, flag, err = arg.Values.Resolve()
		if err != nil {
			flag = cobra.ShellCompDirectiveError
			return
		}
	}

	return
}
