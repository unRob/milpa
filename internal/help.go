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
	"bytes"
	// Embed requires an import so the compiler knows what's up. Golint requires a comment. Gotta please em both.
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

//go:embed help.md
var helpTemplate string

func addBackticks(str []byte) []byte {
	return bytes.ReplaceAll(str, []byte("﹅"), []byte("`"))
}

func readDoc(query []string) ([]byte, error) {
	if len(MilpaPath) == 0 {
		return nil, fmt.Errorf("no MILPA_PATH set on the environment")
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("requesting docs help")
	}

	queryString := strings.Join(query, "/")

	for _, path := range MilpaPath {
		candidate := path + "/docs/" + queryString
		logrus.Debugf("looking for doc named %s", candidate)
		_, err := os.Lstat(candidate + ".md")
		if err == nil {
			return ioutil.ReadFile(candidate + ".md")
		}

		if _, err := os.Lstat(candidate + "/index.md"); err == nil {
			return ioutil.ReadFile(candidate + "/index.md")
		}

		if _, err := os.Stat(candidate); err == nil {
			return []byte{}, BadArguments{fmt.Sprintf("Missing topic for %s", strings.Join(query, " "))}
		}
	}

	return nil, fmt.Errorf("doc not found")
}

var Docs *Command = &Command{
	Summary:     "Dislplays docs on TOPIC",
	Description: "Shows markdown-formatted documentation from milpa repos. See `milpa help docs milpa repo docs` for more information on how to write your own.",
	Arguments: Arguments{
		Argument{
			Name:        "topic",
			Description: "The topic to show docs for",
			Variadic:    true,
			Required:    true,
		},
	},
	Meta: Meta{
		Path: os.Getenv("MILPA_ROOT") + "/milpa/docs",
		Name: []string{"help", "docs"},
		Repo: os.Getenv("MILPA_ROOT"),
		Kind: "internal",
	},
	helpFunc: func(printLinks bool) string {
		topics, err := findDocs([]string{}, "", false)
		if err != nil {
			return ""
		}
		topicList := []string{}
		for _, topic := range topics {
			if printLinks {
				topic = fmt.Sprintf("[%s](%s)", topic, topic)
			}
			topicList = append(topicList, "- "+topic)
		}

		return `## Available topics:

` + strings.Join(topicList, "\n")
	},
}

var DocsCommand *cobra.Command = &cobra.Command{
	Use:   "docs [TOPIC]",
	Short: Docs.Summary,
	Long:  Docs.Description,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		logrus.Debugf("looking for docs given %v and %s", args, toComplete)
		docs, err := findDocs(args, toComplete, false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return docs, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return BadArguments{Msg: "Missing doc topic to display"}
		}

		contents, err := readDoc(args)
		if err != nil {
			switch err.(type) {
			case BadArguments:
				return err
			}
			return NotFound{Msg: "Unknown doc: " + err.Error()}
		}

		titleExp := regexp.MustCompile("^title: (.+)")
		frontmatterSep := []byte("---\n")
		if len(contents) > 3 && string(contents[0:4]) == string(frontmatterSep) {
			// strip out frontmatter
			parts := bytes.SplitN(contents, frontmatterSep, 3)
			title := titleExp.FindString(string(parts[1]))
			if title != "" {
				title = strings.TrimPrefix(title, "title: ")
			} else {
				title = strings.Join(args, " ")
			}
			contents = bytes.Join([][]byte{[]byte("# " + title + "\n"), parts[2]}, []byte("\n"))
		}

		doc, err := renderMD(contents, c)
		if err != nil {
			return err
		}

		if _, err := c.OutOrStderr().Write(doc); err != nil {
			return err
		}
		os.Exit(42)

		return nil
	},
	Annotations: map[string]string{
		"MilpaDocs": "true",
	},
}

var HelpCommand *cobra.Command = &cobra.Command{
	Use:   "help [command]",
	Short: "Display usage information on any **COMMAND...**",
	Long:  `Help provides the valid arguments and options for any command known to milpa. By default, ﹅milpa help﹅ will query the environment variable ﹅COLORFGBG﹅ to decide which style to use when rendering help, except if ﹅MILPA_HELP_STYLE﹅ is set. Valid styles are: **light**, **dark**, and **auto**.`,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		cmd, _, e := c.Root().Find(args)
		if e != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if cmd == nil {
			// Root help command.
			cmd = c.Root()
		}
		for _, subCmd := range cmd.Commands() {
			if subCmd.IsAvailableCommand() || subCmd.Name() == "help" {
				if strings.HasPrefix(subCmd.Name(), toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
				}
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(c *cobra.Command, args []string) {
		cmd, _, e := c.Root().Find(args)
		if cmd == nil || e != nil || (len(args) > 0 && cmd != nil && cmd.Name() != args[len(args)-1]) {
			if cmd == nil {
				err := c.Root().Help()
				if err != nil {
					logrus.Error(err)
					os.Exit(70)
				}
				logrus.Errorf("Unknown help topic %s", args)
				os.Exit(127)
			} else {
				err := cmd.Help()

				if err != nil {
					logrus.Error(err)
					os.Exit(70)
				}

				if len(args) > 1 {
					logrus.Errorf("Unknown help topic %s for %s", args[1], args[0])
				} else {
					logrus.Errorf("Unknown help topic %s for milpa", args[0])
				}
				os.Exit(127)
			}
		} else {
			cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cobra.CheckErr(cmd.Help())
		}

		os.Exit(42)
	},
}

type combinedCommand struct {
	Spec          *Command
	Command       *cobra.Command
	GlobalOptions Options
	HTMLOutput    bool
}

func (cmd *Command) HasAdditionalHelp() bool {
	return cmd.helpFunc != nil
}

func (cmd *Command) AdditionalHelp(printLinks bool) *string {
	if cmd.helpFunc != nil {
		str := cmd.helpFunc(printLinks)
		return &str
	}
	return nil
}

func (cmd *Command) ShowHelp(cc *cobra.Command, args []string) {
	tmpl := template.New("help").Funcs(template.FuncMap{
		"trim":       strings.TrimSpace,
		"toUpper":    strings.ToUpper,
		"trimSuffix": strings.TrimSuffix,
	})
	var err error
	if tmpl, err = tmpl.Parse(helpTemplate); err != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	c := &combinedCommand{
		Spec:          cmd,
		Command:       cc,
		GlobalOptions: Root.Options,
		HTMLOutput:    os.Getenv("MILPA_PLAIN_HELP") == "enabled",
	}
	err = tmpl.Execute(&buf, c)
	if err != nil {
		panic(err)
	}

	content, err := renderMD(buf.Bytes(), cc)
	if err != nil {
		panic(err)
	}

	_, err = cc.OutOrStderr().Write(content)
	if err != nil {
		panic(err)
	}
}

func renderMD(content []byte, cc *cobra.Command) ([]byte, error) {
	content = addBackticks(content)

	if os.Getenv("MILPA_PLAIN_HELP") == "enabled" {
		return content, nil
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return content, err
	}

	var styleFunc glamour.TermRendererOption

	ok, err := cc.Flags().GetBool("no-color")
	if err == nil && ok {
		styleFunc = glamour.WithStandardStyle("notty")
	} else {
		style := os.Getenv("MILPA_HELP_STYLE")
		switch style {
		case "dark":
			styleFunc = glamour.WithStandardStyle("dark")
		case "light":
			styleFunc = glamour.WithStandardStyle("light")
		default:
			styleFunc = glamour.WithStandardStyle("auto")
		}
	}

	renderer, err := glamour.NewTermRenderer(
		styleFunc,
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return content, err
	}

	return renderer.RenderBytes(content)
}
