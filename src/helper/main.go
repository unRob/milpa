package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cmds "github.com/unrob/milpa/command"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func handleError(err error, cmd *cmds.Command, stage string) {
	if err == nil {
		return
	}
	var printHelp = false
	var printError = true
	var exitCode int = 1

	logrus.Debugf("handling error during %s", stage)

	switch v := err.(type) {
	case cmds.BadArguments:
		printHelp = true
		exitCode = 40
	case cmds.NotFound:
		logrus.Debugf("not found: %w", v)
		// if len(v.Group) > 0 {
		logrus.Debugf("looking for valid subcommands for %s", v.Group)
		options := []string{}
		for option, description := range cmds.FindSubCommands(v.Group) {
			options = append(options, fmt.Sprintf("  %s - %s", option, description))
		}
		err = fmt.Errorf("%w. Available sub-commands are: \n%s", err, strings.Join(options, "\n"))
		// }
		exitCode = 44
	default:
		if errors.Is(err, pflag.ErrHelp) {
			printHelp = true
			printError = false
			exitCode = 42
		}
	}

	if printHelp {
		logrus.Debugf("bad args")
		help, err := cmd.Help("markdown")
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(string(help))
	}

	if printError {
		logrus.Error(err)
	}
	os.Exit(exitCode)
}

func main() {
	if os.Getenv("MILPA_VERBOSE") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})

	cmd, remainingArgs, err := cmds.Find(os.Args)
	handleError(err, cmd, "finding command")

	finalArgs, err := cmd.ParseArgs(remainingArgs)
	handleError(err, cmd, "parsing arguments")

	str, err := cmd.ToEval(finalArgs)
	handleError(err, cmd, "printing environment")
	fmt.Println(str)
	if logrus.GetLevel() == logrus.DebugLevel {
		fmt.Println("export MILPA_VERBOSE=1")
	}
}
