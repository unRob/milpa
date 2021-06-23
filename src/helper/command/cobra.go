package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (cmd *Command) ToCobra() (*cobra.Command, error) {
	cc := &cobra.Command{
		Use:               cmd.Meta.Name[len(cmd.Meta.Name)-1],
		Short:             cmd.Summary,
		Long:              cmd.Description,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
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

				if arg.Set != nil {
					values, err := arg.Set.Resolve()
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

	expectedArgLen := len(cmd.Arguments)
	if expectedArgLen > 0 {
		cc.ValidArgsFunction = func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			argsCompleted := len(args)

			values := []string{}
			directive := cobra.ShellCompDirectiveDefault
			// logrus.Infof("argsCompleted: %d, expected: %d", argsCompleted, expectedArgLen)
			if argsCompleted < expectedArgLen {
				// el usuario pide completar un arg que aun esperamos
				arg := cmd.Arguments[argsCompleted]
				if arg.Set != nil {
					values, _ = arg.Set.Resolve()
				} else {
					directive = cobra.ShellCompDirectiveError
				}
			}

			if toComplete != "" {
				filtered := []string{}
				for _, value := range values {
					if strings.HasPrefix(value, toComplete) {
						filtered = append(filtered, value)
					}
				}
				values = filtered
			}

			return values, directive
		}
	}

	err := cmd.CreateFlagSet()
	if err != nil {
		return cc, err
	}

	cc.Flags().AddFlagSet(cmd.runtimeFlags)

	cc.Flags().BoolP("verbose", "v", false, "Log verbose output to stderr")
	cc.Flags().BoolP("help", "h", false, "Display help")

	return cc, nil
}

func RootCommand(commands []*Command) (*cobra.Command, error) {
	root := &cobra.Command{Use: "milpa"}
	root.Flags().AddFlagSet(RootFlagset())

	root.AddCommand(&cobra.Command{
		Use:               "__generate_completions [bash|zsh|fish]",
		Short:             "Outputs a shell-specific script for autocompletions. See milpa help itself shell install-autocomplete",
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		Args:              cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			switch args[0] {
			case "bash":
				err = root.GenBashCompletion(os.Stdout)
			case "zsh":
				err = root.GenZshCompletion(os.Stdout)
			case "fish":
				err = root.GenFishCompletion(os.Stdout, true)
			}
			return
		},
	})

	for _, cmd := range commands {
		leaf, err := cmd.ToCobra()
		if err != nil {
			return nil, err
		}

		container := root
		for idx, cp := range cmd.Meta.Name {
			if idx == len(cmd.Meta.Name)-1 {
				container.AddCommand(leaf)
				break
			}

			query := []string{cp}
			if cc, _, err := root.Find(query); err == nil && cc != container {
				container = cc
			} else {
				cc := &cobra.Command{
					Use:               cp,
					Short:             fmt.Sprintf("%s sub-commands", strings.Join(query, " ")),
					DisableAutoGenTag: true,
					SilenceUsage:      true,
					Args: func(cmd *cobra.Command, args []string) error {
						if len(args) == 0 {
							return NotFound{"No subcommand provided", []string{}}
						}
						if err := cobra.OnlyValidArgs(cmd, args); err != nil {
							return err
						}
						return nil
					},
					RunE: func(cc *cobra.Command, args []string) error {
						if len(args) > 0 {
							fmt.Printf("Error: Unknown subcommand %s\n", args[len(args)-1])
						} else {
							fmt.Println("Error: No subcommand provided")
						}
						err := cc.Help()
						if err != nil {
							logrus.Error(err)
							return err
						}

						return NotFound{"No subcommand provided", []string{}}
					},
				}
				container.AddCommand(cc)
				container = cc
			}
		}
	}

	return root, nil
}
