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
	"fmt"
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
	if err != nil {
		// see man sysexits || grep "#define EX" /usr/include/sysexits.h
		switch err.(type) {
		case internal.BadArguments, internal.ConfigError:
			// 64 bad arguments
			// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
			fmt.Printf("error: %s\n", err)
			showHelp(cmd)

			os.Exit(64)
		case internal.NotFound:
			// 127 command not found
			// https://tldp.org/LDP/abs/html/exitcodes.html
			os.Exit(127)
		default:
			if strings.HasPrefix(err.Error(), "unknown command") {
				os.Exit(127)
			} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
				showHelp(cmd)
				os.Exit(64)
			}
		}

		logrus.Errorf("Unknown error: %s", err)
		os.Exit(2)
	}
}

func main() {
	if os.Getenv("MILPA_DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})

	subcommands, err := internal.FindAllSubCommands(os.Args[1] != "__doctor")
	if err != nil {
		logrus.Fatal(err)
	}

	root, err := internal.RootCommand(subcommands, version)
	// root.SilenceUsage = true
	if err != nil {
		logrus.Errorf("failed to get cobra command for %s: %v", os.Args, err)
		os.Exit(64)
	}

	helpFunc := root.HelpFunc()
	args := []string{}
	initialArgs := []string{"milpa"}
	helpRequested := false

	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		exitCode := 42
		if !helpRequested && cmd.HasAvailableSubCommands() {
			exitCode = 127
		}
		helpFunc(cmd, args)
		os.Exit(exitCode)
	})

	if len(os.Args) > 1 && os.Args[1] != "__complete" {
		for _, arg := range os.Args[1:] {
			if arg == "--help" || arg == "-h" {
				initialArgs = append(initialArgs, "help")
			} else {
				args = append(args, arg)
			}
		}
	} else {
		args = os.Args[1:]
	}

	os.Args = append(initialArgs, args...) //nolint:gocritic

	if len(os.Args) >= 2 && os.Args[1] == "help" {
		helpRequested = true
	}

	handleError(root.ExecuteC())
}
