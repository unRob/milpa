package command

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (cmd *Command) ToCobra() (*cobra.Command, error) {
	cc := &cobra.Command{
		Use:   cmd.Meta.Name[len(cmd.Meta.Name)-1],
		Short: cmd.Summary,
		Long:  cmd.Description,
		Run: func(cc *cobra.Command, args []string) {
			cmd.ToEval(args)
		},
	}

	expectedArgLen := len(cmd.Arguments)
	if expectedArgLen > 0 {
		cc.ValidArgsFunction = func(cc *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			argsCompleted := len(args)
			// last := cmd.Arguments[expectedArgLen-1]
			values := []string{}
			directive := cobra.ShellCompDirectiveDefault
			logrus.Infof("argsCompleted: %d, expected: %d", argsCompleted, expectedArgLen)
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
	// cmd.runtimeFlags.
	cc.Flags().AddFlagSet(cmd.runtimeFlags)

	cc.Flags().BoolP("verbose", "v", false, "Log verbose output to stderr")
	cc.Flags().BoolP("help", "h", false, "Display help for a command")

	return cc, nil
}

func RootCommand(commands []*Command) (*cobra.Command, error) {
	root := &cobra.Command{Use: "milpa"}
	root.Flags().AddFlagSet(RootFlagset())
	//root.()

	for _, cmd := range commands {
		leaf, err := cmd.ToCobra()
		if err != nil {
			return nil, err
		}
		// fmt.Printf("adding command %s\n", cmd.Meta.Name)

		container := root
		for idx, cp := range cmd.Meta.Name {
			if idx == len(cmd.Meta.Name)-1 {
				// fmt.Printf("setting command %s on %s\n", cmd.Meta.Name, container.Name())
				container.AddCommand(leaf)
				// names := []string{}
				// for _, n := range container.Commands() {
				// 	names = append(names, n.Name())
				// }
				// fmt.Printf("subcommands for %s: %v\n", container.Name(), names)
				break
			}

			query := []string{cp}
			// fmt.Printf("searching for container %s in %s\n", query, container.Name())
			if cc, _, err := root.Find(query); err == nil && cc != container {
				container = cc
				// fmt.Printf("found existing container for %s, %v\n", query, cc.Name())
			} else {
				cc := &cobra.Command{
					Use:   cp,
					Short: fmt.Sprintf("%s sub-commands", strings.Join(query, " ")),
					Args:  cobra.MinimumNArgs(1),
					// Run:   func(cc *cobra.Command, args []string) {},
				}
				// fmt.Printf("creating container for %s in %s\n", cp, container.Name())
				container.AddCommand(cc)
				container = cc
			}
			// os.Exit(0)
		}
	}

	return root, nil
}
