// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCommand = &cobra.Command{
	Use:               "__generate_completions [bash|zsh|fish]",
	Short:             "Outputs a shell-specific script for autocompletions. See milpa help itself shell install-autocomplete",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
		return
	},
}
