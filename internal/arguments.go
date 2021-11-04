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
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
	runtime "github.com/unrob/milpa/internal/runtime"
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
type Arguments []Argument

// ToEnv writes shell variables to dst.
func (args *Arguments) ToEnv(dst *[]string, actual []string) error {
	for idx, arg := range *args {
		envName := fmt.Sprintf("%s%s", _c.OutputPrefixArg, strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_")))

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

		if arg.Validates() && runtime.ValidationEnabled() {
			values, _, err := arg.Resolve()
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

// Validate runs validation on provided arguments.
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

		if arg.Validates() && runtime.ValidationEnabled() {
			logrus.Debugf("Validating argument %s", arg.Name)
			values, _, err := arg.Resolve()
			if err != nil {
				return err
			}

			if !contains(values, current) {
				return BadArguments{fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(values, ", "))}
			}
		}
	}

	return nil
}

// CompletionFunction is called by cobra when asked to complete arguments.
func (args *Arguments) CompletionFunction() func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	self := *args
	expectedArgLen := len(self)
	hasVariadicArg := expectedArgLen > 0 && self[len(self)-1].Variadic
	if expectedArgLen > 0 {
		return func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			argsCompleted := len(args)

			values := []string{}
			directive := cobra.ShellCompDirectiveDefault
			if argsCompleted < expectedArgLen || hasVariadicArg {
				var arg Argument
				if hasVariadicArg && argsCompleted >= expectedArgLen {
					// completing a variadic argument
					arg = self[len(self)-1]
				} else {
					// completing regular argument (maybe variadic!)
					arg = self[argsCompleted]
				}
				if arg.providesAutocomplete() {
					values, directive, _ = arg.Resolve()
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
	Values *ValueSource `yaml:"values" validate:"omitempty"`
}

// Validates tells if the user-supplied value needs validation.
func (arg *Argument) Validates() bool {
	return arg.Values != nil && arg.Values.Validates()
}

// providesAutocomplete tells if this option provides autocomplete values.
func (arg *Argument) providesAutocomplete() bool {
	return arg.Values != nil
}

// ToDesc prints out the description of an argument for help and docs.
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
