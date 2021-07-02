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
)

var version = "beta"

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != "help" {
		err := cmd.Help()
		if err != nil {
			os.Exit(70)
		}
	}
}

func handleError(cmd *cobra.Command, err error) {
	if err == nil {
		ok, err := cmd.Flags().GetBool("help")
		if cmd.Name() == "help" || err == nil && ok {
			os.Exit(42)
		}

		os.Exit(0)
	}

	// see man sysexits || grep "#define EX" /usr/include/sysexits.h
	switch err.(type) {
	case internal.BadArguments, internal.ConfigError:
		// 64 bad arguments
		// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(64)
	case internal.NotFound:
		// 127 command not found
		// https://tldp.org/LDP/abs/html/exitcodes.html
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(127)
	default:
		if strings.HasPrefix(err.Error(), "unknown command") {
			showHelp(cmd)
			os.Exit(127)
		} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			showHelp(cmd)
			logrus.Error(err)
			os.Exit(64)
		}
	}
	logrus.Errorf("Unknown error: %s", err)
	os.Exit(2)

}

func main() {
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceColors:            os.Getenv("NO_COLOR") == "",
	})

	isDoctor := false
	if len(os.Args) >= 2 {
		isDoctor = os.Args[1] == "__doctor" || (len(os.Args) > 2 && (os.Args[1] == "itself" && os.Args[2] == "doctor"))
	}

	logrus.Debugf("doctor mode enabled: %v", isDoctor)

	subcommands, err := internal.FindAllSubCommands(!isDoctor)
	if err != nil && !isDoctor {
		logrus.Fatal(err)
	}

	root, err := internal.RootCommand(subcommands, version)
	if err != nil {
		logrus.Errorf("failed to get cobra command for %s: %v", os.Args, err)
		os.Exit(64)
	}

	initialArgs := []string{"milpa"}
	os.Args = append(initialArgs, os.Args[1:]...) //nolint:gocritic

	handleError(root.ExecuteC())
}
