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
package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/runtime"
)

var version = "beta"

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != _c.HelpCommandName {
		err := cmd.Help()
		if err != nil {
			os.Exit(_c.ExitStatusProgrammerError)
		}
	}
}

func handleError(cmd *cobra.Command, err error) {
	if err == nil {
		ok, err := cmd.Flags().GetBool(_c.HelpCommandName)
		if cmd.Name() == _c.HelpCommandName || err == nil && ok {
			os.Exit(_c.ExitStatusRenderHelp)
		}

		os.Exit(_c.ExitStatusOk)
	}

	switch err.(type) {
	case internal.BadArguments, internal.ConfigError:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusUsage)
	case internal.NotFound:
		// 127 command not found
		// https://tldp.org/LDP/abs/html/exitcodes.html
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusNotFound)
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

func main() {
	if os.Getenv(_c.EnvVarDebug) != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceColors:            os.Getenv(_c.EnvVarMilpaUnstyled) == "",
	})

	isDoctor := runtime.DoctorModeEnabled()
	logrus.Debugf("doctor mode enabled: %v", isDoctor)

	subcommands, err := internal.FindAllSubCommands(!isDoctor)
	if err != nil && !isDoctor {
		logrus.Fatal(err)
	}

	root := internal.RootCommand(subcommands, version)

	initialArgs := []string{"milpa"}
	os.Args = append(initialArgs, os.Args[1:]...) //nolint:gocritic

	handleError(root.ExecuteC())
}
