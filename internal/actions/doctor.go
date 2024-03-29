// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"git.rob.mx/nidito/chinampa/pkg/tree"
	"github.com/fatih/color"
	"github.com/unrob/milpa/internal/bootstrap"
	mcmd "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
)

var docLog = logger.Sub("itself doctor")

func DoctorModeEnabled() bool {
	if len(os.Args) < 3 {
		return false
	}
	first := os.Args[1]
	second := os.Args[2]

	return first == "itself" && second == "doctor"
}

var Doctor = &command.Command{
	Path:        []string{"itself", "doctor"},
	Summary:     "Validates all commands found on the `MILPA_PATH`",
	Description: `This command will run checks on all known commands, parsing specs and validating their values.`,
	Options: command.Options{
		"summary": {
			Type:        command.ValueTypeBoolean,
			Description: "Only print errors, if any",
		},
	},
	Action: func(cmd *command.Command) (err error) {
		bold := color.New(color.Bold)
		warn := color.New(color.FgYellow)
		fail := color.New(color.FgRed)
		success := color.New(color.FgGreen)
		failedOverall := false
		failures := map[string]uint8{}
		var out = cmd.Cobra.OutOrStdout()

		summarize := cmd.Options["summary"].ToValue().(bool)

		var milpaRoot string
		if mp := os.Getenv(_c.EnvVarMilpaRoot); mp != "" {
			milpaRoot = strings.Join(strings.Split(mp, ":"), "\n")
		} else {
			milpaRoot = warn.Sprint("empty")
		}
		bold.Fprintf(out, "%s is: %s\n", _c.EnvVarMilpaRoot, milpaRoot)

		var milpaPath string
		bold.Fprintf(out, "%s is: ", _c.EnvVarMilpaPath)
		if mp := os.Getenv(_c.EnvVarMilpaPath); mp != "" {
			milpaPath = "\n" + strings.Join(bootstrap.MilpaPath, "\n")
		} else {
			milpaPath = warn.Sprint("empty")
		}
		fmt.Fprintf(out, "%s\n", milpaPath)
		fmt.Fprintln(out, "")
		bold.Fprintf(out, "Runnable commands:\n")

		for _, cmd := range tree.CommandList() {
			if cmd.Hidden {
				continue
			}
			docLog.Debugf("Validating %s", cmd.FullName())

			message := ""

			hasFailures := false
			report := map[string]int{}
			if meta, ok := cmd.Meta.(mcmd.Meta); ok {
				// fmt.Println("hasmeta")
				parsingErrors := meta.ParsingErrors()
				if len(parsingErrors) > 0 {
					hasFailures = true

					for _, err := range parsingErrors {
						failures[cmd.FullName()]++
						message += fail.Sprintf("  - %s\n", err)
					}
				} else {
					report = cmd.Validate()
				}
			} else {
				report = cmd.Validate()
			}
			for property, status := range report {
				formatter := success
				if status == 1 {
					hasFailures = true
					failures[cmd.FullName()]++
					formatter = fail
				} else if status == 2 {
					formatter = warn
				}

				message += formatter.Sprintf("  - %s\n", property)
			}
			prefix := "✅"
			if hasFailures {
				failedOverall = true
				prefix = "❌"
			}

			fmt.Println(bold.Sprintf("%s %s", prefix, cmd.FullName()), "—", cmd.Path)
			if !summarize || hasFailures {
				if message != "" {
					fmt.Fprintln(out, message)
				}
				fmt.Fprintln(out, "-----------")
			}
		}

		if failedOverall {
			failureReport := []string{}
			for cmd, count := range failures {
				plural := ""
				if count > 1 {
					plural = "s"
				}
				failureReport = append(failureReport, fmt.Sprintf("%s - %d issue%s", cmd, count, plural))
			}

			return fmt.Errorf("your milpa could use some help with the following commands:\n%s", strings.Join(failureReport, "\n"))
		}

		return nil
	},
}
