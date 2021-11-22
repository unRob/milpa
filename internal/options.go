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
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
)

// Options is a map of name to Option.
type Options map[string]*Option

func (opts *Options) AllKnown() map[string]string {
	col := map[string]string{}
	for name, opt := range *opts {
		col[name] = opt.ToString(false)
	}
	return col
}

// ToEnv writes shell variables to dst.
func (opts *Options) ToEnv(command *Command, dst *[]string) {
	command.cc.Flags().VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if name == _c.HelpCommandName {
			return
		}
		envName := ""
		value := f.Value.String()

		if cname, ok := _c.EnvFlagNames[name]; ok {
			if value == "false" {
				return
			}
			envName = cname
			value = shellescape.Quote(value)
		} else {
			envName = fmt.Sprintf("%s%s", _c.OutputPrefixOpt, strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
			oopts := *opts
			opt := oopts[name]
			value = opt.ToString(true)

			if value == "false" && opt.Type == ValueTypeBoolean {
				// makes dealing with false flags in shell easier
				value = ""
			}
		}
		*dst = append(*dst, fmt.Sprintf("export %s=%s", envName, value))
	})
}

func (opts *Options) Parse(supplied *pflag.FlagSet) {
	logrus.Debugf("Parsing supplied flags, %v", supplied)
	for name, opt := range *opts {
		switch opt.Type {
		case ValueTypeBoolean:
			if val, err := supplied.GetBool(name); err == nil {
				opt.provided = val
				continue
			}
		default:
			opt.Type = ValueTypeString
			if val, err := supplied.GetString(name); err == nil {
				opt.provided = val
				continue
			}
		}
	}
}

func (opts *Options) AreValid(cmd *Command) error {
	for name, opt := range *opts {
		if err := opt.Validate(name, cmd); err != nil {
			return err
		}
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
	Repeated    bool         `yaml:"repeated" validate:"omitempty"`
	provided    interface{}
}

func (opt *Option) IsKnown() bool {
	return opt.provided != nil
}

func (opt *Option) ToString(asShell bool) string {
	value := opt.Default
	if opt.IsKnown() {
		value = opt.provided
	}

	stringValue := ""
	switch opt.Type {
	case "bool":
		stringValue = strconv.FormatBool(value.(bool))
	case "string":
		stringValue = value.(string)
	}

	if asShell {
		stringValue = shellescape.Quote(stringValue)
	}

	return stringValue
}

func (opt *Option) Validate(name string, cmd *Command) error {
	if !opt.Validates() {
		return nil
	}

	current := opt.ToString(false) // nolint:ifshort

	if current == "" {
		return nil
	}

	validValues, _, err := opt.Resolve(cmd)
	if err != nil {
		return err
	}

	if !contains(validValues, current) {
		return BadArguments{fmt.Sprintf("%s is not a valid value for option <%s>. Valid options are: %s", current, name, strings.Join(validValues, ", "))}
	}

	return nil
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
func (opt *Option) Resolve(command *Command) (values []string, flag cobra.ShellCompDirective, err error) {
	if opt.Values != nil {
		return opt.Values.Resolve(command)
	}

	return
}

// CompletionFunction is called by cobra when asked to complete an option.
func (opt *Option) CompletionFunction(command *Command) func(cmd *cobra.Command, args []string, toComplete string) (values []string, flag cobra.ShellCompDirective) {
	self := *opt
	return func(cmd *cobra.Command, args []string, toComplete string) (values []string, flag cobra.ShellCompDirective) {
		if !opt.providesAutocomplete() {
			flag = cobra.ShellCompDirectiveNoFileComp
			return
		}

		command.Options.Parse(cmd.Flags())

		var err error
		values, flag, err = self.Resolve(command)
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

}
