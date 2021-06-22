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
		exitCode = 64
		printHelp = true
	case cmds.NotExecutable:
		exitCode = 126
		logrus.Debugf("not executable: %w", v)
	case cmds.NotFound:
		exitCode = 127
		logrus.Debugf("not found: %w", v)
		options := []string{}
		for option, description := range cmds.FindSubCommandDescriptions(v.Group) {
			options = append(options, fmt.Sprintf("  %s - %s", option, description))
		}
		err = fmt.Errorf("%w. Available sub-commands are: \n%s", err, strings.Join(options, "\n"))
	default:
		if errors.Is(err, pflag.ErrHelp) {
			printHelp = true
			printError = false
			exitCode = 42
		}
	}

	if printHelp {
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

func findCommand() {
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

func generateCompletionCommand() {
	cmd, err := cmds.RootCommand([]*cmds.Command{})
	if err != nil {
		logrus.Fatal(err)
	}
	switch os.Args[0] {
	case "bash":
		cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	}
}

func autocompleteCommand() {
	subcommands, err := cmds.FindAllSubCommands()
	if err != nil {
		logrus.Fatal(err)
	}

	root, err := cmds.RootCommand(subcommands)
	if err != nil {
		logrus.Debugf("failed to get cobra command for %s", os.Args)
	}

	args := []string{"milpa", "__complete"}
	if len(os.Args) == 0 {
		os.Args = []string{""}
	}
	os.Args = append(args, os.Args...)
	root.Execute()
}

func main() {
	if os.Getenv("MILPA_DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})

	cmd := ""
	if len(os.Args) < 2 {
		logrus.Fatal("Available helper commands: autocomplete, find")
	}

	cmd = os.Args[1]
	os.Args = os.Args[2:]

	switch cmd {
	case "find":
		findCommand()
	case "__complete":
		autocompleteCommand()
	case "completion":
		generateCompletionCommand()
	default:
		logrus.Errorf("Unknown helper command: %s", cmd)
	}

}
