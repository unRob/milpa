// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"github.com/alessio/shellescape"
	"github.com/spf13/pflag"
	"github.com/unrob/milpa/internal/bootstrap"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/util"
)

// ToEnv writes shell variables to dst.
func ArgumentsToEnv(cmd *command.Command, dst *[]string, prefix string) []string {
	all := []string{}
	for _, arg := range cmd.Arguments {
		envName := fmt.Sprintf("%s%s", _c.OutputPrefixArg, arg.EnvName())

		if arg.Variadic {
			vals := arg.ToValue().([]string)

			ret := []string{}
			for _, v := range vals {
				ret = append(ret, shellescape.Quote(v))
			}
			all = append(all, ret...)
			*dst = append(*dst, fmt.Sprintf("declare -a %s=(%s)", envName, strings.Join(ret, " ")))
		} else {
			*dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, envName, shellescape.Quote(arg.ToString())))
			// all = append(all, arg.ToString())
			// *dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, envName, arg.ToString()))
		}
	}
	return all
}

func ArgumentsToSlice(cmd *command.Command) []string {
	all := []string{}
	for _, arg := range cmd.Arguments {
		if arg.Variadic {
			vals := arg.ToValue().([]string)
			all = append(all, vals...)
		} else {
			all = append(all, arg.ToString())
		}
	}
	return all
}

// FlagNames are flags also available as environment variables.
var flagNames = map[string]string{
	"no-color":        env.NoColor,
	"color":           env.ForceColor,
	"silent":          env.Silent,
	"verbose":         env.Verbose,
	"skip-validation": env.ValidationDisabled,
}

func envValue(opts command.Options, f *pflag.Flag) (*string, *string) {
	name := f.Name
	if name == _c.HelpCommandName {
		return nil, nil
	}
	envName := ""
	value := f.Value.String()

	if cname, ok := flagNames[name]; ok {
		if value == "false" {
			return nil, nil
		}
		return &cname, &value
	}

	envName = fmt.Sprintf("%s%s", _c.OutputPrefixOpt, strings.ToUpper(strings.ReplaceAll(name, "-", "_")))

	if opt := opts[name]; opt != nil {
		if opt.Repeated {
			temp := []string{}
			for _, v := range opt.ToValue().([]string) {
				temp = append(temp, shellescape.Quote(v))
			}
			value = fmt.Sprintf("( %s )", strings.Join(temp, " "))
		} else {
			value = opt.ToString()
		}
	}

	if value == "false" && f.Value.Type() == "bool" {
		// makes dealing with false flags in shell easier
		value = ""
	}

	return &envName, &value
}

// ToEnv writes shell variables to dst.
func OptionsToEnv(cmd *command.Command, dst *[]string, prefix string) {
	cmd.FlagSet().VisitAll(func(f *pflag.Flag) {
		envName, value := envValue(cmd.Options, f)
		if envName != nil && value != nil {
			if opt := cmd.Options[f.Name]; opt != nil && opt.Repeated {
				*dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, *envName, *value))
			} else {
				*dst = append(*dst, fmt.Sprintf("%s%s=%s", prefix, *envName, shellescape.Quote(*value)))
			}
		}
	})
}

func OptionsEnvMap(cmd *command.Command, dst *map[string]string) {
	cmd.Cobra.Flags().VisitAll(func(f *pflag.Flag) {
		envName, value := envValue(cmd.Options, f)
		if envName != nil && value != nil {
			(*dst)[*envName] = *value
		}
	})
}

func EnvironmentMap(cmd *command.Command) map[string]string {
	meta := cmd.Meta.(Meta)
	return map[string]string{
		_c.OutputCommandName: cmd.FullName(),
		_c.OutputCommandKind: string(meta.Kind),
		_c.OutputCommandRepo: meta.Repo,
		_c.OutputCommandPath: meta.Path,
	}
}

func ToEval(cmd *command.Command) string {
	output := []string{}
	for name, value := range util.EnvironmentMap(bootstrap.MilpaPath, bootstrap.MilpaRoot) {
		output = append(output, fmt.Sprintf("export %s=%s", name, shellescape.Quote(value)))
	}

	for name, value := range EnvironmentMap(cmd) {
		output = append(output, fmt.Sprintf("export %s=%s", name, shellescape.Quote(value)))
	}

	OptionsToEnv(cmd, &output, "export ")
	args := ArgumentsToEnv(cmd, &output, "export ")

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	output = append(output, "set -- "+strings.Join(args, " "))

	return strings.Join(output, "\n")
}

func Env(cmd *command.Command, seed []string) []string {
	for name, value := range util.EnvironmentMap(bootstrap.MilpaPath, bootstrap.MilpaRoot) {
		seed = append(seed, fmt.Sprintf("%s=%s", name, shellescape.Quote(value)))
	}

	for name, value := range EnvironmentMap(cmd) {
		seed = append(seed, fmt.Sprintf("%s=%s", name, value))
	}

	OptionsToEnv(cmd, &seed, "")
	ArgumentsToEnv(cmd, &seed, "")

	return seed
}
