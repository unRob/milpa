// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package errors

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	_c "github.com/unrob/milpa/internal/constants"
)

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != _c.HelpCommandName {
		err := cmd.Help()
		if err != nil {
			os.Exit(_c.ExitStatusProgrammerError)
		}
	}
}

func HandleCobraExit(cmd *cobra.Command, err error) {
	if err == nil {
		ok, err := cmd.Flags().GetBool(_c.HelpCommandName)
		if cmd.Name() == _c.HelpCommandName || err == nil && ok {
			os.Exit(_c.ExitStatusRenderHelp)
		}

		os.Exit(_c.ExitStatusOk)
	}

	switch tErr := err.(type) {
	case SubCommandExit:
		logrus.Debugf("Sub-command failed with: %s", err.Error())
		os.Exit(tErr.ExitCode)
	case BadArguments:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusUsage)
	case NotFound:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusNotFound)
	case ConfigError:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusConfigError)
	case EnvironmentError:
		logrus.Error(err)
		os.Exit(_c.ExitStatusConfigError)
	default:
		if strings.HasPrefix(err.Error(), "unknown command") {
			showHelp(cmd)
			os.Exit(_c.ExitStatusNotFound)
		} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			showHelp(cmd)
			logrus.Error(err)
			os.Exit(_c.ExitStatusUsage)
		}
	}

	logrus.Errorf("Unknown error: %s", err)
	os.Exit(2)
}
