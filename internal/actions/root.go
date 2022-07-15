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
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/registry"
	"github.com/unrob/milpa/internal/runtime"
)

var rootcc = &cobra.Command{
	Use: "milpa [--silent|-v|--verbose] [--no-color] [-h|-help] [--version]",
	Annotations: map[string]string{
		_c.ContextKeyRuntimeIndex: "milpa",
	},
	Short: root.Summary,
	Long: `milpa runs commands from .milpa folders

` + root.Description,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	SilenceErrors:     true,
	ValidArgs:         []string{""},
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.OnlyValidArgs(cmd, args)
		if err != nil {

			suggestions := []string{}
			bold := color.New(color.Bold)
			for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
				suggestions = append(suggestions, bold.Sprint(l))
			}
			errMessage := fmt.Sprintf("Unknown subcommand %s", bold.Sprint(strings.Join(args, " ")))
			if len(suggestions) > 0 {
				errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
			}
			return errors.NotFound{Msg: errMessage, Group: []string{}}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if ok, err := cmd.Flags().GetBool("version"); err == nil && ok {
				vc, _, err := cmd.Root().Find([]string{versionCommand.Name()})

				if err != nil {
					return err
				}
				return vc.RunE(vc, []string{})
			}
			return errors.NotFound{Msg: "No subcommand provided", Group: []string{}}
		}

		return nil
	},
}

var root = &command.Command{
	Summary: "Runs commands found in " + _c.RepoRoot + " folders",
	Description: `﹅milpa﹅ is a command-line tool to care for one's own garden of scripts, its name comes from "milpa", an agricultural method that combines multiple crops in close proximity. You and your team write scripts and a little spec for each command -use bash, or any other language-, and ﹅milpa﹅ provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.

  See [﹅milpa help docs milpa﹅](/.milpa/docs/milpa/index.md) for more information about ﹅milpa﹅`,
	Meta: command.Meta{
		Path: _c.EnvVarMilpaRoot + "/" + _c.Milpa,
		Name: []string{_c.Milpa},
		Repo: _c.EnvVarMilpaRoot,
		Kind: command.KindRoot,
	},
	Options: command.Options{
		_c.HelpCommandName: &command.Option{
			ShortName:   "h",
			Type:        "bool",
			Description: "Display help for any command",
		},
		"verbose": &command.Option{
			ShortName:   "v",
			Type:        "bool",
			Default:     runtime.VerboseEnabled(),
			Description: "Log verbose output to stderr",
		},
		"no-color": &command.Option{
			Type:        "bool",
			Description: "Print to stderr without any formatting codes",
		},
		"silent": &command.Option{
			Type:        "bool",
			Description: "Silence non-error logging",
		},
		"skip-validation": &command.Option{
			Type:        "bool",
			Description: "Do not validate any arguments or options",
		},
	},
}

func RootCommand(version string) *cobra.Command {
	rootcc.Annotations["version"] = version
	rootFlagset := pflag.NewFlagSet("compa", pflag.ContinueOnError)
	rootFlagset.BoolP("verbose", "v", runtime.VerboseEnabled(), "Log verbose output to stderr")
	rootFlagset.Bool("silent", false, "Do not print any logs to stderr")
	rootFlagset.BoolP("help", "h", false, "Display help for any command")
	rootFlagset.Bool("no-color", !runtime.ColorEnabled(), "Do not print any formatting codes")
	rootFlagset.Bool("skip-validation", false, "Do not validate any arguments or options")
	rootFlagset.Usage = func() {}
	rootFlagset.SortFlags = false

	rootcc.PersistentFlags().AddFlagSet(rootFlagset)
	rootcc.Flags().Bool("version", false, "Display the version of milpa")

	rootcc.CompletionOptions.DisableDefaultCmd = true
	rootcc.AddCommand(versionCommand)
	rootcc.AddCommand(completionCommand)
	rootcc.AddCommand(generateDocumentationCommand)
	rootcc.AddCommand(doctorCommand)
	rootcc.AddCommand(fetchCommand)
	rootcc.AddCommand(introspectCommand)
	introspectCommand.Flags().Int32("depth", 15, "")
	introspectCommand.Flags().String("format", "json", "")
	introspectCommand.Flags().String("template", "{{ indent . }}{{ .Name }} - {{ .Summary }}\n", "")
	rootcc.SetHelpCommand(helpCommand)
	helpCommand.AddCommand(docsCommand)
	docsCommand.SetHelpFunc(docs.HelpRenderer(root.Options))
	rootcc.SetHelpFunc(root.HelpRenderer(root.Options))

	registry.SetRoot(rootcc, root)
	return rootcc
}
