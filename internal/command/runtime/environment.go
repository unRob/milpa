// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"github.com/alessio/shellescape"
	"github.com/spf13/pflag"
	"github.com/unrob/milpa/internal/bootstrap"
	"github.com/unrob/milpa/internal/command/meta"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/util"
)

// ToEnv writes shell variables to dst.
func ArgumentsToEnv(cmd *command.Command, dst *[]string) []string {
	m := cmd.Meta.(meta.Meta)
	all := []string{}
	for _, arg := range cmd.Arguments {
		envVarName := envVarName(arg.EnvName(), _c.OutputPrefixArg)
		value := envVarPair(*envVarName, arg, m)
		*dst = append(*dst, *value)

		if arg.Repeats() {
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

// ToEnv writes shell variables to dst.
func OptionsToEnv(cmd *command.Command, dst *[]string) {
	m := cmd.Meta.(meta.Meta)
	cmd.FlagSet().VisitAll(func(f *pflag.Flag) {
		name := f.Name
		if name == _c.HelpCommandName {
			return
		}
		envVar := envVarName(name, _c.OutputPrefixOpt)

		// check if part of global flags
		if _, ok := flagNames[name]; ok {
			value := f.Value.String()
			if value == "false" {
				return
			}
			*dst = append(*dst, *envVarValue(m, *envVar, value))
			return
		}

		if opt := cmd.Options[name]; opt != nil {
			if f.Value.Type() == "bool" && opt.ToString() == "false" {
				return
			}

			value := envVarPair(*envVar, opt, m)
			*dst = append(*dst, *value)
		}
	})
}

func EnvironmentMap(cmd *command.Command) map[string]string {
	meta := cmd.Meta.(meta.Meta)
	return map[string]string{
		_c.OutputCommandName: cmd.FullName(),
		_c.OutputCommandKind: string(meta.Kind),
		_c.OutputCommandRepo: meta.Repo,
		_c.OutputCommandPath: meta.Path,
	}
}

// ToEval returns a sequence of commands to be interpreted by a shell.
func ToEval(cmd *command.Command) string {
	m := cmd.Meta.(meta.Meta)
	output := []string{}
	for name, value := range util.EnvironmentMap(bootstrap.MilpaPath, bootstrap.MilpaRoot) {
		output = append(output, *envVarValue(m, name, value))
	}

	for name, value := range EnvironmentMap(cmd) {
		output = append(output, *envVarValue(m, name, value))
	}

	OptionsToEnv(cmd, &output)
	args := ArgumentsToEnv(cmd, &output)

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	output = append(output, "set -- "+strings.Join(args, " "))

	res := strings.Join(output, "\n")
	return res
}

func Env(cmd *command.Command, seed []string) ([]string, []string) {
	m := cmd.Meta.(meta.Meta)
	for name, value := range util.EnvironmentMap(bootstrap.MilpaPath, bootstrap.MilpaRoot) {
		seed = append(seed, *envVarValue(m, name, value))
	}

	for name, value := range EnvironmentMap(cmd) {
		seed = append(seed, *envVarValue(m, name, value))
	}

	OptionsToEnv(cmd, &seed)
	args := ArgumentsToEnv(cmd, &seed)

	return seed, args
}

func BaseEnv(m meta.Meta) []string {
	itself, err := os.Executable()
	if err != nil {
		log.Debugf("could not determine milpa's executable path: %s", err)
	}

	env := []string{
		_c.EnvVarMilpaRoot + "=" + bootstrap.MilpaRoot,
		_c.OutputCommandPath + "=" + m.Path,
		"MILPA=" + itself,
	}
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if strings.HasPrefix(parts[0], "MILPA_COMMAND_") ||
			strings.HasPrefix(parts[0], "MILPA_ARG_") ||
			strings.HasPrefix(parts[0], "MILPA_OPT_") ||
			parts[0] == _c.EnvVarMilpaRoot ||
			parts[0] == _c.OutputCommandPath {
			continue
		}
		env = append(env, kv)
	}
	return env
}
