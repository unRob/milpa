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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
)

var helpCommand = &cobra.Command{
	Use:   _c.HelpCommandName + " [command]",
	Short: "Display usage information on any **COMMAND...**",
	Long:  `Help provides the valid arguments and options for any command known to milpa. By default, ﹅milpa help﹅ will query the environment variable ﹅COLORFGBG﹅ to decide which style to use when rendering help, except if ﹅MILPA_HELP_STYLE﹅ is set. Valid styles are: **light**, **dark**, and **auto**.`,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		cmd, _, e := c.Root().Find(args)
		if e != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if cmd == nil {
			// Root help command.
			cmd = c.Root()
		}
		for _, subCmd := range cmd.Commands() {
			if subCmd.IsAvailableCommand() || subCmd.Name() == _c.HelpCommandName {
				if strings.HasPrefix(subCmd.Name(), toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
				}
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(c *cobra.Command, args []string) {
		cmd, _, e := c.Root().Find(args)
		if cmd == nil || e != nil || (len(args) > 0 && cmd != nil && cmd.Name() != args[len(args)-1]) {
			if cmd == nil {
				err := c.Root().Help()
				if err != nil {
					logrus.Error(err)
					os.Exit(_c.ExitStatusProgrammerError)
				}
				logrus.Errorf("Unknown help topic %s", args)
				os.Exit(_c.ExitStatusNotFound)
			} else {
				err := cmd.Help()

				if err != nil {
					logrus.Error(err)
					os.Exit(_c.ExitStatusProgrammerError)
				}

				if len(args) > 1 {
					logrus.Errorf("Unknown help topic %s for %s", args[1], args[0])
				} else {
					logrus.Errorf("Unknown help topic %s for milpa", args[0])
				}
				os.Exit(_c.ExitStatusNotFound)
			}
		} else {
			cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cobra.CheckErr(cmd.Help())
		}

		os.Exit(_c.ExitStatusRenderHelp)
	},
}
