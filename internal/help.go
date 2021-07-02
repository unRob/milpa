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
	"unicode"

	"github.com/charmbracelet/glamour"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func sighGolang(str string) string {
	return strings.ReplaceAll(str, "", "`")
}

func findDocs(query []string, needle string) ([]string, error) {
	results := []string{}
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
			if needle == "" || strings.HasPrefix(name, needle) {
				results = append(results, name)
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

	queryString := strings.Join(query, "/") + ".md"
	logrus.Debugf("looking for doc %s in %s", queryString, MilpaPath)

	for _, path := range MilpaPath {
		candidate := path + "/.milpa/docs/" + queryString
		logrus.Debugf("looking for doc in %s", candidate)
		if _, err := os.Stat(candidate); err == nil {
			return ioutil.ReadFile(candidate)
		}
	}

	return nil, fmt.Errorf("doc not found")
}

var DocsCommand *cobra.Command = &cobra.Command{
	Use:   "docs [topic]",
	Short: "Docs on specific topics",
	Long:  `docs shows formatted documentation from milpa repos`,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		logrus.Debugf("looking for docs given %v and %s", args, toComplete)
		docs, err := findDocs(args, toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return docs, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(c *cobra.Command, args []string) error {
		contents, err := readDoc(args)
		if err != nil {
			return NotFound{Msg: "Unknown doc: " + err.Error()}
		}

		width, _, err := terminal.GetSize(0)
		if err != nil {
			return err
		}

		logrus.Debug(len(contents))
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
		if cmd == nil || e != nil {
			c.Printf("Unknown help topic %#q\n", args)
			cobra.CheckErr(c.Root().Usage())
		} else {
			cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cobra.CheckErr(cmd.Help())
		}

		// return nil
	},
}

// func (cmd *Command) Help

type combinedCommand struct {
	Spec    *Command
	Command *cobra.Command
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

func (cmd *Command) ShowHelp(cc *cobra.Command, args []string) {
	tmpl := template.New("help").Funcs(template.FuncMap{
		"trim":                    strings.TrimSpace,
		"trimRightSpace":          trimRightSpace,
		"trimTrailingWhitespaces": trimRightSpace,
		// "appendIfNotPresent":      appendIfNotPresent,
		"rpad": rpad,
		// "gt":                      Gt,
		// "eq":                      Eq,
		"toUpper": strings.ToUpper,
	})
	var err error
	if tmpl, err = tmpl.Parse(HelpTemplate); err != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	c := &combinedCommand{
		Spec:    cmd,
		Command: cc,
	}
	err = tmpl.Execute(&buf, c)
	if err != nil {
		panic(err)
	}

	width, _, err := terminal.GetSize(0)
	if err != nil {
		panic(err)
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		panic(err)
	}

	help, err := renderer.Render(buf.String())
	if err != nil {
		panic(err)
	}
	fmt.Println(help)
	os.Exit(42)
}

var HelpTemplate = sighGolang(`# milpa {{ .Spec.FullName }}

{{ .Command.Short }}

## Usage:

  {{ .Command.UseLine }}{{if .Command.HasAvailableSubCommands}} subcommand{{end}}

## Description:

{{ .Spec.Description }}

{{ if .Command.HasAvailableSubCommands -}}
## Subcommands:

{{- range .Command.Commands -}}
{{- if (or .IsAvailableCommand (eq .Name "help")) -}}
- {{rpad .Name .NamePadding }} - {{.Short}}
{{ end -}}
{{- end -}}
{{- end -}}

{{- if .Spec.Arguments -}}
## Arguments:

{{ range .Spec.Arguments -}}
- {{ if .Required}}**{{ end }}${{ .Name | toUpper }}{{ if .Variadic}}...{{ end }}{{ if .Required }}**{{ end }} - {{ .Description }}
{{ end -}}
{{- end -}}

{{- if .Command.HasAvailableLocalFlags}}
## Options:

{{ range $name, $opt := .Spec.Options -}}
- --{{ $name }} (_{{$opt.Type}}_): {{$opt.Description}}.{{ if $opt.Default }} Default: _{{ $opt.Default }}_.{{ end }}
{{ end -}}
{{- end -}}

{{- if .Command.HasAvailableInheritedFlags -}}
## Global Options:

{{.Command.InheritedFlags.FlagUsages | trimTrailingWhitespaces }}

{{end}}
`)
