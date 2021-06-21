package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var MILPA_PATH []string = strings.Split(os.Getenv("MILPA_PATH"), ":")

func Find(args []string) (*Command, []string, error) {
	var finalError error
	flag := pflag.NewFlagSet("helper", pflag.ContinueOnError)
	flag.BoolP("verbose", "v", false, "Log verbose output to stderr")
	flag.BoolP("help", "h", false, "Display help for a command")
	flag.Usage = func() {}
	flag.SortFlags = false

	err := flag.Parse(args)

	subcommand := args[1:]
	logrus.Debugf("original args %s", subcommand)

	if ok, err := flag.GetBool("verbose"); err == nil && ok {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Verbose output enabled")
		for idx, arg := range subcommand {
			if arg == "--verbose" || arg == "-v" {
				subcommand = append(subcommand[0:idx], subcommand[idx+1:]...)
			}
		}
	}

	if err != nil {
		if !strings.HasPrefix(err.Error(), "unknown flag: ") {
			logrus.Error(err)
			os.Exit(2)
		}
	}

	if ok, err := flag.GetBool("help"); err == nil && ok {
		finalError = pflag.ErrHelp
		logrus.Debug(finalError)
		for idx, v := range subcommand {
			if v == "--help" {
				subcommand = append(subcommand[:idx], subcommand[idx+1:]...)
			}
		}
	}

	if len(subcommand) == 0 {
		return nil, args, NotFound{Msg: "No command provided"}
	}

	logrus.Debugf("found args %s", subcommand)

	for _, pkg := range MILPA_PATH {
		for i := range subcommand {
			if i == len(subcommand) {
				logrus.Debug("done with args, %d")
				break
			}
			query := subcommand[0 : len(subcommand)-i]
			logrus.Debugf("looking for %s", query)
			searchPath := fmt.Sprintf("%s/.milpa/commands/%s", pkg, strings.Join(query, "/"))
			commandPath := fmt.Sprintf("%s.sh", searchPath)
			_, err := os.Stat(commandPath)

			kind := "source"
			if err != nil {
				commandPath = searchPath
				statInfo, err := os.Stat(commandPath)
				if err != nil {
					continue
				}

				if statInfo.IsDir() {
					var msg error
					if i == 0 {
						msg = NotFound{
							Msg:   fmt.Sprintf("missing sub-command for %s", strings.Join(query, " ")),
							Group: query,
						}
					} else {
						msg = NotFound{
							Msg:   fmt.Sprintf("found command group for <%s>, but no sub-command named <%s>", strings.Join(query, " "), subcommand[len(subcommand)-i]),
							Group: query,
						}
					}
					return nil, args, msg
				}

				if statInfo.Mode()&0100 == 0 {
					return nil, args, fmt.Errorf("found %s but it's not executable", commandPath)
				}
				kind = "exec"
			}

			logrus.Debugf("Found script: %s", commandPath)
			cmd, err := New(commandPath, fmt.Sprintf("%s.yaml", searchPath), pkg, kind)
			if err != nil {
				logrus.Debugf("failed to construct command for %v", query)
				return nil, args, err
			}

			return cmd, subcommand[len(query):], finalError
		}
	}

	return nil, args, NotFound{Msg: fmt.Sprintf("No command found named %s", subcommand[0])}
}

func FindSubCommands(query []string) map[string]string {
	results := map[string]string{}
	for _, path := range MILPA_PATH {
		queryBase := fmt.Sprintf("%s/.milpa/commands/%s", path, strings.Join(query, "/"))

		matches, err := filepath.Glob(fmt.Sprintf("%s/*", queryBase))
		logrus.Debugf("found matches %v", matches)
		if err == nil {
			for _, match := range matches {
				if !(strings.HasSuffix(match, ".sh") || filepath.Ext(match) == "") {
					logrus.Debugf("ignoring %s", match)
					continue
				}
				logrus.Debugf("reading %s", match)

				fileInfo, err := os.Stat(match)
				if err != nil {
					continue
				}

				if !fileInfo.IsDir() {
					pc := strings.SplitN(match, "/.milpa/commands/", 2)
					pkg := pc[0]
					path := pc[1]

					spec := fmt.Sprintf("%s/.milpa/commands/%s.yaml", pkg, strings.Replace(path, ".sh", "", 1))
					cmd, err := New(match, spec, pkg, "tbd")
					if err != nil {
						continue
					}
					name := strings.Replace(fileInfo.Name(), ".sh", "", 1)
					results[name] = cmd.Summary
				} else {
					// assume everything here is a sub-command group, figure out a fix later
					results[fileInfo.Name()] = fmt.Sprintf("%s sub-commands", fileInfo.Name())
				}
			}
		}
	}

	return results
}
