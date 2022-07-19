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
