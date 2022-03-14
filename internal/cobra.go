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
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
	runtime "github.com/unrob/milpa/internal/runtime"
)

var rootFlagset *pflag.FlagSet

func RootFlagset() *pflag.FlagSet {
	if rootFlagset == nil {

		verboseDefault := runtime.VerboseEnabled()
		noColor := !runtime.ColorEnabled()

		rootFlagset = pflag.NewFlagSet("compa", pflag.ContinueOnError)
		rootFlagset.BoolP("verbose", "v", verboseDefault, "Log verbose output to stderr")
		rootFlagset.Bool("silent", false, "Do not print any logs to stderr")
		rootFlagset.BoolP("help", "h", false, "Display help for any command")
		rootFlagset.Bool("no-color", noColor, "Do not print any formatting codes")
		rootFlagset.Bool("skip-validation", false, "Do not validate any arguments or options")
		rootFlagset.Usage = func() {}
		rootFlagset.SortFlags = false
	}

	return rootFlagset
}

func (cmd *Command) Run(cc *cobra.Command, args []string) error {
	cmd.Arguments.Parse(args)
	cmd.Options.Parse(cc.Flags())
	if runtime.ValidationEnabled() {
		if err := cmd.Arguments.AreValid(cmd); err != nil {
			return err
		}
	}
	if runtime.ValidationEnabled() {
		if err := cmd.Options.AreValid(cmd); err != nil {
			return err
		}
	}

	env := cmd.ToEval(args)

	if os.Getenv(_c.EnvVarCompaOut) != "" {
		return os.WriteFile(os.Getenv(_c.EnvVarCompaOut), []byte(env), 0600)
	}

	fmt.Println(env)
	return nil
}

func (cmd *Command) ToCobra() *cobra.Command {
	if cmd.cc != nil {
		return cmd.cc
	}

	localName := cmd.Meta.Name[len(cmd.Meta.Name)-1]
	useSpec := []string{localName, "[options]"}
	for _, arg := range cmd.Arguments {
		useSpec = append(useSpec, arg.ToDesc())
	}

	cc := &cobra.Command{
		Use:               strings.Join(useSpec, " "),
		Short:             cmd.Summary,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args: func(cc *cobra.Command, supplied []string) error {
			if runtime.ValidationEnabled() {
				cmd.Arguments.Parse(supplied)
				return cmd.Arguments.AreValid(cmd)
			}
			return nil
		},
		RunE: cmd.Run,
	}

	cc.SetFlagErrorFunc(func(c *cobra.Command, e error) error {
		return BadArguments{e.Error()}
	})

	cc.ValidArgsFunction = cmd.Arguments.CompletionFunction(cmd)

	cmd.CreateFlagSet()
	cc.Flags().AddFlagSet(cmd.runtimeFlags)

	for name, opt := range cmd.Options {
		if err := cc.RegisterFlagCompletionFunc(name, opt.CompletionFunction(cmd)); err != nil {
			logrus.Errorf("Failed setting up autocompletion for option <%s> of command <%s>", name, cmd.FullName())
		}
	}

	cc.SetHelpFunc(cmd.ShowHelp)
	cmd.cc = cc
	return cc
}

func RootCommand(commands []*Command, version string) *cobra.Command {
	root := &cobra.Command{
		Use:         "milpa [--silent|-v|--verbose] [--no-color] [-h|-help] [--version]",
		Annotations: map[string]string{"version": version},
		Short:       Root.Summary,
		Long: `milpa runs commands from .milpa folders

` + Root.Description,
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
				return NotFound{errMessage, []string{}}
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
				return NotFound{"No subcommand provided", []string{}}
			}

			return nil
		},
	}
	root.PersistentFlags().AddFlagSet(RootFlagset())
	root.Flags().Bool("version", false, "Display the version of milpa")

	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(versionCommand)
	root.AddCommand(completionCommand)
	root.AddCommand(generateDocumentationCommand)
	root.AddCommand(doctorForCommands(commands))
	root.AddCommand(fetchCommand)
	root.SetHelpCommand(HelpCommand)
	HelpCommand.AddCommand(DocsCommand)
	DocsCommand.SetHelpFunc(Docs.ShowHelp)
	root.SetHelpFunc(Root.ShowHelp)

	populateRoot(root, commands)
	return root
}

func populateRoot(root *cobra.Command, commands []*Command) {
	for _, cmd := range commands {
		leaf := cmd.ToCobra()

		container := root
		for idx, cp := range cmd.Meta.Name {
			if idx == len(cmd.Meta.Name)-1 {
				// logrus.Debugf("adding command %s to %s", leaf.Name(), container.Name())
				container.AddCommand(leaf)
				break
			}

			query := []string{cp}
			if cc, _, err := container.Find(query); err == nil && cc != container {
				logrus.Debugf("found %s in %s", query, cc.Name())
				container = cc
			} else {
				logrus.Debugf("creating %s in %s", query, container.Name())
				cc := &cobra.Command{
					Use:                        cp,
					Short:                      fmt.Sprintf("%s subcommands", strings.Join(query, " ")),
					DisableAutoGenTag:          true,
					SuggestionsMinimumDistance: 2,
					SilenceUsage:               true,
					SilenceErrors:              true,
					Args: func(cmd *cobra.Command, args []string) error {
						err := cobra.OnlyValidArgs(cmd, args)
						if err != nil {

							suggestions := []string{}
							bold := color.New(color.Bold)
							for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
								suggestions = append(suggestions, bold.Sprint(l))
							}
							last := len(args) - 1
							parent := cmd.CommandPath()
							errMessage := fmt.Sprintf("Unknown subcommand %s of known command %s", bold.Sprint(args[last]), bold.Sprint(parent))
							if len(suggestions) > 0 {
								errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
							}
							return NotFound{errMessage, []string{}}
						}
						return nil
					},
					ValidArgs: []string{""},
					RunE: func(cc *cobra.Command, args []string) error {
						if len(args) == 0 {
							return NotFound{"No subcommand provided", []string{}}
						}
						os.Exit(_c.ExitStatusNotFound)
						return nil
					},
				}
				groupParent := &Command{
					Summary:     fmt.Sprintf("%s subcommands", strings.Join(query, " ")),
					Description: fmt.Sprintf("Runs subcommands within %s", strings.Join(query, " ")),
					Arguments:   Arguments{},
					Options:     Options{},
					Meta: Meta{
						Name: cmd.Meta.Name[0 : idx+1],
					},
				}
				cc.SetHelpFunc(groupParent.ShowHelp)
				container.AddCommand(cc)
				container = cc
			}
		}
	}
}
