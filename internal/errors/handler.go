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
	case BadArguments, ConfigError:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(_c.ExitStatusUsage)
	case NotFound:
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
