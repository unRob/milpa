package errors

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ConfigError struct {
	Err    error
	Config string
}

type EnvironmentError struct {
	Err error
}

func (err ConfigError) Error() string {
	return fmt.Sprintf("Invalid configuration %s: %v", err.Config, err.Err)
}

func (err EnvironmentError) Error() string {
	return fmt.Sprintf("Invalid MILPA_ environment: %v", err.Err)
}

func showHelp(cmd *cobra.Command) {
	if cmd.Name() != "help" {
		err := cmd.Help()
		if err != nil {
			os.Exit(statuscode.ProgrammerError)
		}
	}
}

func HandleCobraExit(cmd *cobra.Command, err error) error {
	if err == nil {
		ok, err := cmd.Flags().GetBool("help")
		if cmd.Name() == "help" || err == nil && ok {
			os.Exit(statuscode.RenderHelp)
		}

		os.Exit(statuscode.Ok)
	}

	switch err.(type) {
	case errors.BadArguments:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(statuscode.Usage)
	case errors.NotFound:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(statuscode.NotFound)
	case ConfigError:
		showHelp(cmd)
		logrus.Error(err)
		os.Exit(statuscode.ConfigError)
	case EnvironmentError:
		logrus.Error(err)
		os.Exit(statuscode.ConfigError)
	default:
		if strings.HasPrefix(err.Error(), "unknown command") {
			showHelp(cmd)
			os.Exit(statuscode.NotFound)
		} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			showHelp(cmd)
			logrus.Error(err)
			os.Exit(statuscode.Usage)
		}
	}

	logrus.Errorf("Unknown error: %s", err)
	os.Exit(2)
	return err
}
