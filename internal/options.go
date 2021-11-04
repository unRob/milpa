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
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
	runtime "github.com/unrob/milpa/internal/runtime"
)

// Options is a map of name to Option.
type Options map[string]*Option

// ToEnv writes shell variables to dst.
func (opts *Options) ToEnv(dst *[]string, flags *pflag.FlagSet) (err error) {
	errors := []string{}
	flags.VisitAll(func(f *pflag.Flag) {
		name := f.Name
		//nolint:goconst
		if name == "help" {
			return
		}
		envName := ""
		value := f.Value.String()

		if cname, ok := _c.EnvFlagNames[name]; ok {
			if value == "false" {
				return
			}
			envName = cname
		} else {
			envName = fmt.Sprintf("%s%s", _c.OutputPrefixOpt, strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
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
			if value != "" && ok && opt.Validates() && runtime.ValidationEnabled() {
				logrus.Debugf("Validating option %s", name)
				values, _, verr := opt.Resolve()
				if verr != nil {
					errors = append(errors, err.Error())
					return
				}

				if !contains(values, value) {
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

// Option represents a command line flag.
type Option struct {
	ShortName   string       `yaml:"short-name"`
	Type        ValueType    `yaml:"type" validate:"omitempty,oneof=string bool"`
	Description string       `yaml:"description" validate:"required"`
	Default     interface{}  `yaml:"default"`
	Values      *ValueSource `yaml:"values" validate:"omitempty"`
}

// Validates tells if the user-supplied value needs validation.
func (opt *Option) Validates() bool {
	return opt.Values != nil && opt.Values.Validates()
}

// providesAutocomplete tells if this option provides autocomplete values.
func (opt *Option) providesAutocomplete() bool {
	return opt.Values != nil
}

// Resolve returns autocomplete values for an option.
func (opt *Option) Resolve() (values []string, flag cobra.ShellCompDirective, err error) {
	if opt.Values != nil {
		return opt.Values.Resolve()
	}

	return
}

// CompletionFunction is called by cobra when asked to complete an option.
func (opt *Option) CompletionFunction(cmd *cobra.Command, args []string, toComplete string) (values []string, flag cobra.ShellCompDirective) {

	if !opt.providesAutocomplete() {
		flag = cobra.ShellCompDirectiveError
		return
	}

	var err error
	values, flag, err = opt.Resolve()
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
	return values, flag
}
