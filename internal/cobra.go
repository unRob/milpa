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

func (cmd *Command) ToCobra() (*cobra.Command, error) {
	localName := cmd.Meta.Name[len(cmd.Meta.Name)-1]
	useSpec := []string{localName, "[flags]"}
	for _, arg := range cmd.Arguments {
		useSpec = append(useSpec, arg.ToDesc())
	}

	// logrus.Debugf("Cobraizing %s", strings.Join(useSpec, " "))
	cc := &cobra.Command{
		Use:               strings.Join(useSpec, " "),
		Short:             cmd.Summary,
		Long:              color.New(color.Bold).Sprintf(strings.Join(useSpec, " ")) + "\n\n" + cmd.Summary + "\n\n" + cmd.Description,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args: func(cc *cobra.Command, args []string) error {
			for idx, arg := range cmd.Arguments {
				argumentProvided := idx < len(args)
				if arg.Required && !argumentProvided {
					return BadArguments{fmt.Sprintf("Missing argument for %s", arg.Name)}
				}

				if !argumentProvided {
					continue
				}
				current := args[idx]

				if arg.Validates() {
					values, err := arg.Resolve()
					if err != nil {
						return err
					}
					found := false
					for _, value := range values {
						if value == current {
							found = true
							break
						}
					}

					if !found {
						return BadArguments{fmt.Sprintf("%s is not a valid value for argument <%s>. Valid options are: %s", current, arg.Name, strings.Join(values, ", "))}
					}
				}
			}

			return nil
		},
		RunE: func(cc *cobra.Command, args []string) error {
			flags := cc.Flags()
			env, err := cmd.ToEval(args, flags)
			if err != nil {
				return err
			}
			fmt.Println(env)
			return nil
		},
	}

	cc.SetFlagErrorFunc(func(c *cobra.Command, e error) error {
		return BadArguments{e.Error()}
	})

	cc.ValidArgsFunction = cmd.Arguments.ToValidationFunction()

	err := cmd.CreateFlagSet()
	if err != nil {
		return cc, err
	}

	cc.Flags().AddFlagSet(cmd.runtimeFlags)

	for name, opt := range cmd.Options {
		if opt.Validates() {
			oCopy := opt
			cc.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				logrus.Infof("registering complete function for %s", name)
				values, err := oCopy.Resolve()
				if err != nil {
					return values, cobra.ShellCompDirectiveError
				}

				if toComplete != "" {
					logrus.Infof("returning filtered results for %s", name)
					filtered := []string{}
					for _, value := range values {
						if strings.HasPrefix(value, toComplete) {
							filtered = append(filtered, value)
						}
					}
					values = filtered
				} else {
					logrus.Infof("returning all results for %s, %v", name, values)
				}

				return values, cobra.ShellCompDirectiveDefault
			})
		}
	}

	return cc, nil
}

func RootCommand(commands []*Command, version string) (*cobra.Command, error) {
	bold := color.New(color.Bold)
	root := &cobra.Command{
		Use:     "milpa subcommand... [--silent|-v|--verbose] [--no-color] [-h|-help]",
		Version: version,
		Short:   "milpa runs other scripts from .milpa folders",
		Long: `milpa runs other scripts from .milpa folders

Milpa, is an agricultural method that combines multiple crops in close proximity. ` + bold.Sprintf("milpa") + ` is a Bash script and tool to care for one's own garden of scripts. You and your team write scripts and a little spec for each command. Use bash, or any other command, and ` + bold.Sprintf("milpa") + ` provides autocompletions, sub-commands and argument parsing+validation for you to focus on your scripts.`,
		DisableAutoGenTag: true,
	}
	root.PersistentFlags().AddFlagSet(RootFlagset())

	root.AddCommand(completionCommand)
	root.AddCommand(doctorForCommands(commands))

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
					// SilenceUsage:      true,
					Args: func(cmd *cobra.Command, args []string) error {
						err := cobra.OnlyValidArgs(cmd, args)
						if err != nil {
							suggestions := []string{}
							bold := color.New(color.Bold)

							for _, l := range strings.Split(err.Error(), "\n") {
								if strings.HasPrefix(l, "\t") {
									suggestions = append(suggestions, bold.Sprint(strings.TrimPrefix(l, "\t")))
								}
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
				container.AddCommand(cc)
				container = cc
			}
		}
	}

	return root, nil
}
