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
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
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

func envValue(opts Options, f *pflag.Flag) (*string, *string) {
	name := f.Name
	if name == _c.HelpCommandName {
		return nil, nil
	}
	envName := ""
	value := f.Value.String()

	if cname, ok := _c.EnvFlagNames[name]; ok {
		if value == "false" {
			return nil, nil
		}
		envName = cname
		value = shellescape.Quote(value)
	} else {
		envName = fmt.Sprintf("%s%s", _c.OutputPrefixOpt, strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
		opt := opts[name]
		value = opt.ToString(true)

		if value == "false" && opt.Type == ValueTypeBoolean {
			// makes dealing with false flags in shell easier
			value = ""
		}
	}

	return &envName, &value
}

// ToEnv writes shell variables to dst.
func (opts *Options) ToEnv(command *Command, dst *[]string, prefix string) {
	command.cc.Flags().VisitAll(func(f *pflag.Flag) {
		envName, value := envValue(*opts, f)
		if envName != nil && value != nil {
			*dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, *envName, *value))
		}
	})
}

func (opts *Options) EnvMap(command *Command, dst *map[string]string) {
	command.cc.Flags().VisitAll(func(f *pflag.Flag) {
		envName, value := envValue(*opts, f)
		if envName != nil && value != nil {
			(*dst)[*envName] = *value
		}
	})
}

func (opts *Options) Parse(supplied *pflag.FlagSet) {
	// logrus.Debugf("Parsing supplied flags, %v", supplied)
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

func (opts *Options) AreValid() error {
	for name, opt := range *opts {
		if err := opt.Validate(name); err != nil {
			return err
		}
	}

	return nil
}

// Option represents a command line flag.
type Option struct {
	ShortName   string       `json:"short-name,omitempty" yaml:"short-name,omitempty"` // nolint:tagliatelle
	Type        ValueType    `json:"type" yaml:"type" validate:"omitempty,oneof=string bool"`
	Description string       `json:"description" yaml:"description" validate:"required"`
	Default     interface{}  `json:"default,omitempty" yaml:"default,omitempty"`
	Values      *ValueSource `json:"values,omitempty" yaml:"values,omitempty" validate:"omitempty"`
	Repeated    bool         `json:"repeated" yaml:"repeated" validate:"omitempty"`
	Command     *Command     `json:"-" yaml:"-" validate:"-"`
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
		if value == nil {
			stringValue = ""
		} else {
			stringValue = strconv.FormatBool(value.(bool))
		}
	case "string":
		if value != nil {
			stringValue = value.(string)
		}
	}

	if asShell {
		stringValue = shellescape.Quote(stringValue)
	}

	return stringValue
}

func (opt *Option) Validate(name string) error {
	if !opt.Validates() {
		return nil
	}

	current := opt.ToString(false) // nolint:ifshort

	if current == "" {
		return nil
	}

	validValues, _, err := opt.Resolve()
	if err != nil {
		return err
	}

	if !contains(validValues, current) {
		return errors.BadArguments{Msg: fmt.Sprintf("%s is not a valid value for option <%s>. Valid options are: %s", current, name, strings.Join(validValues, ", "))}
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
func (opt *Option) Resolve() (values []string, flag cobra.ShellCompDirective, err error) {
	if opt.Values != nil {
		if opt.Values.Command == nil {
			opt.Values.Command = opt.Command
		}
		return opt.Values.Resolve()
	}

	return
}

// CompletionFunction is called by cobra when asked to complete an option.
func (opt *Option) CompletionFunction(cmd *cobra.Command, args []string, toComplete string) (values []string, flag cobra.ShellCompDirective) {
	if !opt.providesAutocomplete() {
		flag = cobra.ShellCompDirectiveNoFileComp
		return
	}

	opt.Command.Arguments.Parse(args)
	opt.Command.Options.Parse(cmd.Flags())

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

	return cobra.AppendActiveHelp(values, opt.Description), flag
}
