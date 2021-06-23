package main

import (
	"os"
	"strings"

	cmds "github.com/unrob/milpa/command"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if os.Getenv("MILPA_DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})

	subcommands, err := cmds.FindAllSubCommands()
	if err != nil {
		logrus.Fatal(err)
	}

	root, err := cmds.RootCommand(subcommands)
	root.SilenceUsage = true
	if err != nil {
		logrus.Debugf("failed to get cobra command for %s", os.Args)
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

	os.Args = append(initialArgs, args...)
	// logrus.Info(os.Args)
	if len(os.Args) >= 2 && os.Args[1] == "help" {
		helpRequested = true
	}

	err = root.Execute()
	if err != nil {
		// see man sysexits || grep "#define EX" /usr/include/sysexits.h
		switch err.(type) {
		case cmds.BadArguments:
			// 64 bad arguments
			// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
			os.Exit(64)
		case cmds.NotFound:
			// 127 command not found
			// https://tldp.org/LDP/abs/html/exitcodes.html
			os.Exit(127)
		default:
			if strings.HasPrefix(err.Error(), "unknown command") {
				os.Exit(127)
			} else if strings.HasPrefix(err.Error(), "unknown flag") || strings.HasPrefix(err.Error(), "unknown shorthand flag") {
				os.Exit(64)
			}
		}

		logrus.Errorf("Unknown error: %s", err)
		os.Exit(2)
	}
}
