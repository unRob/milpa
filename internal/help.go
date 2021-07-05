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
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func addBackticks(str []byte) []byte {
	return bytes.ReplaceAll(str, []byte("﹅"), []byte("`"))
}

func findDocs(query []string, needle string) ([]string, error) {
	results := []string{}
	found := map[string]bool{}
	if len(MilpaPath) == 0 {
		return results, fmt.Errorf("no MILPA_PATH set on the environment")
	}

	logrus.Debugf("looking for docs in %s", MilpaPath)
	queryString := ""
	if len(query) > 0 {
		queryString = strings.Join(query, "/") + "/"
	}

	for _, path := range MilpaPath {
		qbase := path + "/.milpa/docs/" + queryString
		q := qbase + "/*"
		logrus.Debugf("looking for docs matching %s", q)
		docs, err := filepath.Glob(q)
		if err != nil {
			return results, err
		}

		for _, doc := range docs {
			name := strings.TrimSuffix(filepath.Base(doc), ".md")
			if _, ok := found[name]; (needle == "" || strings.HasPrefix(name, needle)) && !ok {
				results = append(results, name)
				found[name] = true
			}
		}

	}

	return results, nil
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
		candidate := path + "/.milpa/docs/" + queryString
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
	Description: "docs shows markdown-formatted documentation from " + os.Getenv("MILPA_NAME") + ` repos. See ` + "`" + os.Getenv("MILPA_NAME") + " help docs milpa repo docs` for more information on how to write your own.",
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
	helpFunc: func() string {
		topics, err := findDocs([]string{}, "")
		if err != nil {
			return ""
		}
		topicList := []string{}
		for _, topic := range topics {
			topicList = append(topicList, "- "+topic)
		}

		return `## Available topics:

` + strings.Join(topicList, "\n")
	},
}

var DocsCommand *cobra.Command = &cobra.Command{
	Use:   "docs [TOPIC]",
	Short: "Displays docs on TOPIC",
	Long:  "docs shows formatted documentation from " + os.Getenv("MILPA_NAME") + " repos",
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		logrus.Debugf("looking for docs given %v and %s", args, toComplete)
		docs, err := findDocs(args, toComplete)
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

		frontmatterSep := []byte("---\n")
		// logrus.Info(string(contents[0:4]) == "---\n")
		if len(contents) > 3 && string(contents[0:4]) == string(frontmatterSep) {
			// strip out frontmatter
			parts := bytes.SplitN(contents, frontmatterSep, 3)
			contents = parts[2]
		}

		contents = bytes.ReplaceAll(contents, []byte("!milpa!"), []byte(os.Getenv("MILPA_NAME")))

		width, _, err := term.GetSize(0)
		if err != nil {
			return err
		}

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithEmoji(),
			glamour.WithWordWrap(width),
		)

		if err != nil {
			return err
		}

		doc, err := renderer.RenderBytes(contents)
		if err != nil {
			return err
		}
		fmt.Println(string(doc))
		os.Exit(42)

		return nil
	},
}

var HelpCommand *cobra.Command = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long:  `Help provides help for any command in the application.`,
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
	Bin           string
}

func (cmd *Command) HasAdditionalHelp() bool {
	return cmd.helpFunc != nil
}

func (cmd *Command) AdditionalHelp() *string {
	if cmd.helpFunc != nil {
		str := cmd.helpFunc()
		return &str
	}
	return nil
}

func (cmd *Command) ShowHelp(cc *cobra.Command, args []string) {
	tmpl := template.New("help").Funcs(template.FuncMap{
		"trim":    strings.TrimSpace,
		"toUpper": strings.ToUpper,
	})
	var err error
	if tmpl, err = tmpl.Parse(HelpTemplate); err != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	c := &combinedCommand{
		Spec:          cmd,
		Command:       cc,
		GlobalOptions: Root.Options,
		Bin:           os.Getenv("MILPA_NAME"),
	}
	err = tmpl.Execute(&buf, c)
	if err != nil {
		panic(err)
	}
	content := addBackticks(buf.Bytes())

	if os.Getenv("MILPA_PLAIN_HELP") == "enabled" {
		_, err = cc.OutOrStderr().Write(content)
		if err != nil {
			panic(err)
		} else {
			return
		}
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	style := glamour.WithAutoStyle
	ok, err := cc.Flags().GetBool("no-color")
	if err == nil && ok {
		style = func() glamour.TermRendererOption { return glamour.WithStyles(glamour.ASCIIStyleConfig) }
	}

	renderer, err := glamour.NewTermRenderer(
		style(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		panic(err)
	}

	help, err := renderer.RenderBytes(content)
	if err != nil {
		panic(err)
	}
	_, err = cc.OutOrStderr().Write(help)
	if err != nil {
		panic(err)
	}
}

var HelpTemplate = `# {{ if not (eq .Spec.Meta.Kind "root") }}{{ .Bin }} {{ end }}{{ .Spec.FullName }}

{{ .Command.Short }}

## Usage:

  ﹅{{ .Command.UseLine }}{{if .Command.HasAvailableSubCommands}} subcommand{{end}}﹅

{{ if .Command.HasAvailableSubCommands -}}
## Subcommands:

{{ range .Command.Commands -}}
{{- if (or .IsAvailableCommand (eq .Name "help")) -}}
- ﹅{{ .Name }}﹅ - {{.Short}}
{{ end -}}
{{- end -}}
{{ else }}
{{- if .Spec.Arguments -}}
## Arguments:

{{ range .Spec.Arguments -}}

- {{ if .Required}}**{{ end }}﹅{{ .Name | toUpper }}{{ if .Variadic}}...{{ end }}﹅{{ if .Required }}**{{ end }} - {{ .Description }}
{{ end -}}
{{- end -}}

## Description:

{{ .Spec.Description }}
{{- if .Spec.HasAdditionalHelp }}
{{ .Spec.AdditionalHelp }}
{{ end -}}

{{ end -}}
{{ if eq .Spec.Meta.Kind "root" }}
## Description:

{{ .Spec.Description }}
{{ end -}}



{{- if .Command.HasAvailableLocalFlags}}
## Options:

{{ range $name, $opt := .Spec.Options -}}
- ﹅--{{ $name }}﹅ (_{{$opt.Type}}_): {{$opt.Description}}.{{ if $opt.Default }} Default: _{{ $opt.Default }}_.{{ end }}
{{ end -}}
{{- end -}}

{{- if .Command.HasAvailableInheritedFlags }}
## Global Options:

{{ range $name, $opt := .GlobalOptions -}}
- ﹅--{{ $name }}﹅ (_{{$opt.Type}}_): {{$opt.Description}}.{{ if $opt.Default }} Default: _{{ $opt.Default }}_.{{ end }}
{{ end -}}
{{end}}`
