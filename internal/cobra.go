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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootFlagset *pflag.FlagSet

func RootFlagset() *pflag.FlagSet {
	if rootFlagset == nil {

		verboseDefault := os.Getenv("MILPA_VERBOSE") != ""
		colorDefault := os.Getenv("NO_COLOR") != ""

		rootFlagset = pflag.NewFlagSet("compa", pflag.ContinueOnError)
		rootFlagset.BoolP("verbose", "v", verboseDefault, "Log verbose output to stderr")
		rootFlagset.BoolP("help", "h", false, "Display help for any command")
		rootFlagset.Bool("no-color", colorDefault, "Do not print any formatting codes")
		rootFlagset.Bool("silent", false, "Do not print any logs to stderr")
		rootFlagset.Usage = func() {}
		rootFlagset.SortFlags = false
	}

	return rootFlagset
}

func (cmd *Command) Run(cc *cobra.Command, args []string) error {
	flags := cc.Flags()
	env, err := cmd.ToEval(args, flags)
	if err != nil {
		return err
	}
	fmt.Println(env)
	return nil
}

func (cmd *Command) ToCobra() (*cobra.Command, error) {
	localName := cmd.Meta.Name[len(cmd.Meta.Name)-1]
	useSpec := []string{localName, "[options]"}
	for _, arg := range cmd.Arguments {
		useSpec = append(useSpec, arg.ToDesc())
	}

	// logrus.Debugf("Cobraizing %s", strings.Join(useSpec, " "))
	cc := &cobra.Command{
		Use:   strings.Join(useSpec, " "),
		Short: cmd.Summary,
		// Long:              color.New(color.Bold).Sprintf(strings.Join(useSpec, " ")) + "\n\n" + cmd.Summary + "\n\n" + cmd.Description,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cmd.Arguments.Validate,
		RunE:              cmd.Run,
	}

	cc.SetFlagErrorFunc(func(c *cobra.Command, e error) error {
		return BadArguments{e.Error()}
	})

	cc.ValidArgsFunction = cmd.Arguments.CompletionFunction()

	err := cmd.CreateFlagSet()
	if err != nil {
		return cc, err
	}

	cc.Flags().AddFlagSet(cmd.runtimeFlags)

	for name, opt := range cmd.Options {
		if opt.Validates() {
			if err := cc.RegisterFlagCompletionFunc(name, opt.ValidationFunction); err != nil {
				return cc, err
			}
		}
	}

	cc.SetHelpFunc(cmd.ShowHelp)

	return cc, nil
}

func RootCommand(commands []*Command, version string) (*cobra.Command, error) {
	root := &cobra.Command{
		Use:     os.Getenv("MILPA_NAME") + " [--silent|-v|--verbose] [--no-color] [-h|-help]",
		Version: version,
		Short:   os.Getenv("MILPA_NAME") + " runs commands from .milpa folders",
		Long: `milpa runs commands from .milpa folders

Milpa, is an agricultural method that combines multiple crops in close proximity. ﹅milpa﹅ is a Bash script and tool to care for one's own garden of scripts. You and your team write scripts and a little spec for each command. Use bash, or any other command, and ﹅milpa﹅ provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.`,
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
				errMessage := fmt.Sprintf("Unknown subcommand %s", args[len(args)-1])
				if len(suggestions) > 0 {
					errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
				}
				return NotFound{errMessage, []string{}}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return NotFound{"No subcommand provided", []string{}}
			}

			return nil
		},
	}
	root.PersistentFlags().AddFlagSet(RootFlagset())

	root.AddCommand(completionCommand)
	root.AddCommand(doctorForCommands(commands))
	root.SetHelpCommand(HelpCommand)
	HelpCommand.AddCommand(DocsCommand)
	DocsCommand.SetHelpFunc(Docs.ShowHelp)
	root.SetHelpFunc(Root.ShowHelp)

	for _, cmd := range commands {
		leaf, err := cmd.ToCobra()
		if err != nil {
			return nil, err
		}

		container := root
		for idx, cp := range cmd.Meta.Name {
			if idx == len(cmd.Meta.Name)-1 {
				// logrus.Debugf("adding command %s to %s", leaf.Name(), container.Name())
				container.AddCommand(leaf)
				break
			}

			query := []string{cp}
			if cc, _, err := container.Find(query); err == nil && cc != container {
				// logrus.Debugf("found %s in %s", query, cc.Name())
				container = cc
			} else {
				// logrus.Debugf("creating %s in %s", query, container.Name())
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
							errMessage := fmt.Sprintf("Unknown subcommand %s", args[len(args)-1])
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
						os.Exit(127)
						return nil
					},
				}
				ccmd := &Command{
					Summary:     fmt.Sprintf("%s subcommands", strings.Join(query, " ")),
					Description: fmt.Sprintf("Runs subcommands within %s", strings.Join(query, " ")),
					Arguments:   Arguments{},
					Options:     Options{},
					Meta: Meta{
						Name: cmd.Meta.Name[0 : idx+1],
					},
				}
				cc.SetHelpFunc(ccmd.ShowHelp)
				container.AddCommand(cc)
				container = cc
			}
		}
	}

	return root, nil
}
