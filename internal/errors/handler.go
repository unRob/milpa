// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package errors

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/spf13/cobra"
)

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != "help" {
		err := cmd.Help()
		if err != nil {
			os.Exit(statuscode.ProgrammerError)
		}
	}
}

func HandleExit(cmd *cobra.Command, err error) error {
	if err == nil {
		ok, err := cmd.Flags().GetBool("help")
		if cmd.Name() == "help" || err == nil && ok {
			os.Exit(statuscode.Ok)
		}

		os.Exit(statuscode.Ok)
	}

	switch e := err.(type) {
	case errors.BadArguments:
		showHelp(cmd)
		logger.Error(err)
		os.Exit(statuscode.Usage)
	case errors.NotFound:
		showHelp(cmd)
		logger.Error(err)
		os.Exit(statuscode.NotFound)
	case SpecError:
		fmt.Println(e.Error())
		os.Exit(statuscode.ProgrammerError)
	case EnvironmentError:
		logger.Error(err)
		os.Exit(statuscode.ConfigError)
	case ProgrammerError:
		logger.Error(err)
		os.Exit(statuscode.ProgrammerError)
	default:
		if strings.HasPrefix(err.Error(), "unknown command") {
			showHelp(cmd)
			os.Exit(statuscode.NotFound)
		} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			showHelp(cmd)
			logger.Error(err)
			os.Exit(statuscode.Usage)
		}
		logger.Errorf("Unknown error: (%#v) %s", err, err)
	}

	os.Exit(2)
	return err
}
