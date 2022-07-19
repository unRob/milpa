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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	command "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/registry"
	"gopkg.in/yaml.v2"
)

var introspectCommand = &cobra.Command{
	Use:               "__inspect [prefix...]",
	Short:             "Inspects the command tree at a given PREFIX",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 1 && args[0] == "" {
			args = []string{}
		}
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
		if err != nil {
			format = "json"
		}

		var serializationFn func(interface{}) ([]byte, error)
		switch format {
		case "yaml":
			serializationFn = yaml.Marshal
		case "json":
			serializationFn = func(t interface{}) ([]byte, error) { return json.MarshalIndent(t, "", "  ") }
		case "text":
			outputTpl, err := cmd.Flags().GetString("template")
			if err != nil {
				outputTpl = "{{ .Name }} - {{ .Summary }}\n"
			}

			tpl := template.Must(template.New("treeItem").Funcs(_c.TemplateFuncs).Parse(outputTpl))
			serializationFn = func(t interface{}) ([]byte, error) {
				tree := t.(*registry.CommandTree)
				var output bytes.Buffer
				err := tree.Traverse(func(cmd *command.Command) error { return tpl.Execute(&output, cmd) })
				return output.Bytes(), err
			}
		case "autocomplete":
			serializationFn = func(interface{}) ([]byte, error) {
				return []byte(strings.Join(registry.ChildrenNames(), "\n") + "\n"), nil
			}
		default:
			return errors.BadArguments{Msg: fmt.Sprintf("Unknown format <%s> for command tree serialization", format)}
		}

		serialized, err := registry.SerializeTree(serializationFn)
		if err != nil {
			return err
		}
		fmt.Print(serialized)

		return nil
	},
}
