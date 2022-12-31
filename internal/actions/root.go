// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

// var rootcc = &cobra.Command{
// 	Use: "milpa [--silent|-v|--verbose] [--[no-]color] [-h|--help] [--version]",
// 	Annotations: map[string]string{
// 		_c.ContextKeyRuntimeIndex: "milpa",
// 	},
// 	Short: Root.Summary,
// 	Long: `milpa runs commands from .milpa folders

// ` + Root.Description,
// 	DisableAutoGenTag: true,
// 	SilenceUsage:      true,
// 	SilenceErrors:     true,
// 	ValidArgs:         []string{""},
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		err := cobra.OnlyValidArgs(cmd, args)
// 		if err != nil {

// 			suggestions := []string{}
// 			bold := color.New(color.Bold)
// 			for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
// 				suggestions = append(suggestions, bold.Sprint(l))
// 			}
// 			errMessage := fmt.Sprintf("Unknown subcommand %s", bold.Sprint(strings.Join(args, " ")))
// 			if len(suggestions) > 0 {
// 				errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
// 			}
// 			return errors.NotFound{Msg: errMessage, Group: []string{}}
// 		}
// 		return nil
// 	},
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if len(args) == 0 {
// 			if ok, err := cmd.Flags().GetBool("version"); err == nil && ok {
// 				vc, _, err := cmd.Root().Find([]string{versionCommand.Name()})

// 				if err != nil {
// 					return err
// 				}
// 				return vc.RunE(vc, []string{})
// 			}
// 			return errors.NotFound{Msg: "No subcommand provided", Group: []string{}}
// 		}

// 		return nil
// 	},
// }

// var Root = &command.Command{
// 	Summary: "Runs commands found in " + _c.RepoRoot + " folders",
// 	Description: `﹅milpa﹅ is a command-line tool to care for one's own garden of scripts, its name comes from an agricultural method that combines multiple crops in close proximity. You and your team write scripts and a little spec for each command -use bash, or any other language-, and ﹅milpa﹅ provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.

//   See [﹅milpa help docs milpa﹅](/.milpa/docs/milpa/index.md) for more information about ﹅milpa﹅`,
// 	Meta: command.Meta{
// 		Path: _c.EnvVarMilpaRoot + "/" + _c.Milpa,
// 		Name: []string{_c.Milpa},
// 		Repo: _c.EnvVarMilpaRoot,
// 		Kind: command.KindRoot,
// 	},
// 	Options: command.Options{
// 		_c.HelpCommandName: &command.Option{
// 			ShortName:   "h",
// 			Type:        "bool",
// 			Description: "Display help for any command",
// 		},
// 		"verbose": &command.Option{
// 			ShortName:   "v",
// 			Type:        "bool",
// 			Default:     runtime.VerboseEnabled(),
// 			Description: "Log verbose output to stderr",
// 		},
// 		"no-color": &command.Option{
// 			Type:        "bool",
// 			Description: "Disable printing of colors to stderr",
// 			Default:     !runtime.ColorEnabled(),
// 		},
// 		"color": &command.Option{
// 			Type:        "bool",
// 			Description: "Always print colors to stderr",
// 			Default:     runtime.ColorEnabled(),
// 		},
// 		"silent": &command.Option{
// 			Type:        "bool",
// 			Description: "Silence non-error logging",
// 		},
// 		"skip-validation": &command.Option{
// 			Type:        "bool",
// 			Description: "Do not validate any arguments or options",
// 		},
// 	},
// }

// func RootCommand(version string) *cobra.Command {
// 	rootcc.Annotations["version"] = version
// 	rootFlagset := pflag.NewFlagSet("compa", pflag.ContinueOnError)
// 	for name, opt := range Root.Options {
// 		def, ok := opt.Default.(bool)
// 		if !ok {
// 			def = false
// 		}

// 		if opt.ShortName != "" {
// 			rootFlagset.BoolP(name, opt.ShortName, def, opt.Description)
// 		} else {
// 			rootFlagset.Bool(name, def, opt.Description)
// 		}
// 	}

// 	rootFlagset.Usage = func() {}
// 	rootFlagset.SortFlags = false
// 	rootcc.PersistentFlags().AddFlagSet(rootFlagset)

// 	rootcc.Flags().Bool("version", false, "Display the version of milpa")

// 	rootcc.CompletionOptions.DisableDefaultCmd = true

// 	rootcc.AddCommand(versionCommand)
// 	rootcc.AddCommand(completionCommand)
// 	rootcc.AddCommand(generateDocumentationCommand)
// 	rootcc.AddCommand(doctorCommand)

// 	doctorCommand.Flags().Bool("summary", false, "")
// 	rootcc.AddCommand(fetchRemoteRepo)
// 	rootcc.AddCommand(introspectCommand)

// 	introspectCommand.Flags().Int32("depth", 15, "")
// 	introspectCommand.Flags().String("format", "json", "")
// 	introspectCommand.Flags().String("template", "{{ indent . }}{{ .Name }} - {{ .Summary }}\n", "")

// 	rootcc.SetHelpCommand(helpCommand)
// 	helpCommand.AddCommand(docsCommand)
// 	docsCommand.SetHelpFunc(docs.HelpRenderer(Root.Options))
// 	rootcc.SetHelpFunc(Root.HelpRenderer(Root.Options))

// 	return rootcc
// }
