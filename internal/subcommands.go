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
package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var completionCommand *cobra.Command = &cobra.Command{
	Use:               "__generate_completions [bash|zsh|fish]",
	Short:             "Outputs a shell-specific script for autocompletions. See milpa help itself shell install-autocomplete",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
		return
	},
}

func doctorForCommands(commands []*Command) *cobra.Command {
	return &cobra.Command{
		Use:               "__doctor",
		Short:             "Outputs information about milpa and known repos. See milpa help itself doctor",
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		RunE: func(_ *cobra.Command, args []string) (err error) {
			bold := color.New(color.Bold)
			warn := color.New(color.FgYellow)
			fail := color.New(color.FgRed)
			success := color.New(color.FgGreen)

			var milpaRoot string
			if mp := os.Getenv("MILPA_ROOT"); mp != "" {
				milpaRoot = strings.Join(strings.Split(mp, ":"), "\n")
			} else {
				milpaRoot = warn.Sprint("empty")
			}
			bold.Printf("MILPA_ROOT is: %s\n", milpaRoot)

			var milpaPath string
			bold.Printf("MILPA_PATH is: ")
			if mp := os.Getenv("MILPA_PATH"); mp != "" {
				milpaPath = "\n" + strings.Join(strings.Split(mp, ":"), "\n")
			} else {
				milpaPath = warn.Sprint("empty")
			}
			fmt.Printf("%s\n", milpaPath)
			fmt.Println("")
			bold.Printf("Runnable commands:\n")

			sort.Sort(ByPath(commands))
			for _, cmd := range commands {
				report := cmd.Validate()
				message := ""

				hasFailures := false
				for property, isValid := range report {
					formatter := success
					if !isValid {
						hasFailures = true
						formatter = fail
					}

					message += formatter.Sprintf("  - %s\n", property)
				}
				prefix := "✅"
				if hasFailures {
					prefix = "❌"
				}

				fmt.Println(bold.Sprintf("%s %s", prefix, cmd.FullName()), "—", cmd.Meta.Path)
				fmt.Println(message)
				fmt.Println("-----------")
			}
			return
		},
	}
}
