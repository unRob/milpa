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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal/registry"
)

var introspectCommand *cobra.Command = &cobra.Command{
	Use:               "__inspect [prefix...]",
	Short:             "Inspects the command tree at a given PREFIX",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		base, remainingArgs, err := cmd.Root().Find(args)

		if err != nil {
			return err
		}

		if len(remainingArgs) > 0 && len(remainingArgs) == len(args) {
			return nil
		}

		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			depth = 20
		}

		logrus.Debugf("looking for commands at %s depth: %d", base.Name(), depth)
		registry.BuildTree(base, depth)
		format, err := cmd.Flags().GetString("format")

		if err == nil {
			switch format {
			case "json":
				json, err := registry.AsJSONTree()
				if err != nil {
					return err
				}
				fmt.Print(json)
			case "autocomplete":
				for _, name := range registry.ChildrenNames() {
					fmt.Println(name)
				}
			}

		}

		return nil
	},
}
