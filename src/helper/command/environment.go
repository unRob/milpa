package command

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func (cmd *Command) ToEval(args []string) (string, error) {
	envVars := []string{
		fmt.Sprintf("export MILPA_COMMAND_NAME=%s", shellescape.Quote(cmd.FullName())),
		fmt.Sprintf("export MILPA_COMMAND_KIND=%s", shellescape.Quote(cmd.Meta.Kind)),
		fmt.Sprintf("export MILPA_COMMAND_PACKAGE=%s", shellescape.Quote(cmd.Meta.Package)),
		fmt.Sprintf("export MILPA_COMMAND_PATH=%s", shellescape.Quote(cmd.Meta.Path)),
	}

	cmd.runtimeFlags.VisitAll(func(f *pflag.Flag) {
		envName := fmt.Sprintf("MILPA_OPT_%s", strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))

		value := f.Value.String()
		switch f.Value.Type() {
		case "bool":
			if val, err := cmd.runtimeFlags.GetBool(f.Name); err == nil && !val {
				value = ""
			} else {
				value = "true"
			}
		default:
			logrus.Debugf("flag %s is a %s", f.Name, f.Value.Type())
		}

		value = shellescape.Quote(value)
		envVars = append(envVars, fmt.Sprintf("export %s=%s", envName, value))
	})

	logrus.Debugf("Printing environment for args: %v", args)

	for idx, arg := range cmd.Arguments {
		envName := fmt.Sprintf("MILPA_ARG_%s", strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_")))

		if idx >= len(args) {
			if arg.Required {
				return "", fmt.Errorf("missing argument: %s", arg.Name)
			}
			logrus.Debugf("Skipping arg parsing for %s", arg.Name)
			envVars = append(envVars, fmt.Sprintf("export %s=%s", envName, ""))
			continue
		}

		var value string
		if arg.Variadic {
			values := []string{}
			for _, va := range args[idx:] {
				values = append(values, shellescape.Quote(va))
			}
			value = fmt.Sprintf("( %s )", strings.Join(values, " "))

		} else {
			value = shellescape.Quote(args[idx])
		}

		if arg.Set != nil {
			values, err := arg.Set.Resolve()
			if err != nil {
				return "", err
			}
			found := false
			for _, validValue := range values {
				if value == validValue {
					found = true
					break
				}
			}

			if !found {
				return "", BadArguments{
					fmt.Sprintf("invalid value for %s: %s. Valid values are %s", arg.Name, value, strings.Join(values, ", ")),
				}
			}
		}

		logrus.Debugf("arg parsing for %s, %d, %s", arg.Name, idx, value)

		envVars = append(envVars, fmt.Sprintf("export %s=%s", envName, value))
	}

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	envVars = append(envVars, "set -- "+strings.Join(args, " "))

	return strings.Join(envVars, "\n"), nil
}
