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
package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/registry"
	"github.com/unrob/milpa/internal/runtime"
)

var doctorCommand = &cobra.Command{
	Use:               "__doctor",
	Short:             "Outputs information about milpa and known repos. See milpa help itself doctor",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		bold := color.New(color.Bold)
		warn := color.New(color.FgYellow)
		fail := color.New(color.FgRed)
		success := color.New(color.FgGreen)
		failedOverall := false
		failures := map[string]uint8{}

		summarize, err := cmd.Flags().GetBool("summary")
		if err != nil {
			summarize = false
		}

		var milpaRoot string
		if mp := os.Getenv(_c.EnvVarMilpaRoot); mp != "" {
			milpaRoot = strings.Join(strings.Split(mp, ":"), "\n")
		} else {
			milpaRoot = warn.Sprint("empty")
		}
		bold.Printf("%s is: %s\n", _c.EnvVarMilpaRoot, milpaRoot)

		var milpaPath string
		bold.Printf("%s is: ", _c.EnvVarMilpaPath)
		if mp := os.Getenv(_c.EnvVarMilpaPath); mp != "" {
			milpaPath = "\n" + strings.Join(runtime.MilpaPath, "\n")
		} else {
			milpaPath = warn.Sprint("empty")
		}
		fmt.Printf("%s\n", milpaPath)
		fmt.Println("")
		bold.Printf("Runnable commands:\n")

		for _, cmd := range registry.CommandList() {
			logrus.Debugf("Validating %s", cmd.FullName())
			report := cmd.Validate()
			message := ""

			hasFailures := false
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

			fmt.Println(bold.Sprintf("%s %s", prefix, cmd.FullName()), "—", cmd.Meta.Path)
			if !summarize || hasFailures {
				if message != "" {
					fmt.Println(message)
				}
				fmt.Println("-----------")
			}
		}

		if failedOverall {
			failureReport := []string{}
			for cmd, count := range failures {
				failureReport = append(failureReport, fmt.Sprintf("%s - %d issues", cmd, count))
			}

			return fmt.Errorf("your milpa could use some help with the following commands:\n%s", strings.Join(failureReport, "\n"))
		}

		return
	},
}
