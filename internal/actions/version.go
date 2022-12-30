// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
)

var versionCommand = &cobra.Command{
	Use:               "__version",
	Short:             "Display the version of milpa",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output := cmd.ErrOrStderr()
		version := cmd.Root().Annotations["version"]
		if cmd.CalledAs() == "" {
			// user asked for --version directly
			output = cmd.OutOrStderr()
			version += "\n"
		}

		_, err := output.Write([]byte(version))
		if err != nil {
			logrus.Errorf("version error: %s", err)
			return err
		}

		os.Exit(_c.ExitStatusRenderHelp)
		return nil
	},
}
