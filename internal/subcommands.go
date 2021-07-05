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
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var completionCommand *cobra.Command = &cobra.Command{
	Use:               "__generate_completions [bash|zsh|fish]",
	Short:             "Outputs a shell-specific script for autocompletions. See milpa help itself shell install-autocomplete",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
		return
	},
}

func doctorForCommands(commands []*Command) *cobra.Command {
	return &cobra.Command{
		Use:               "__doctor",
		Short:             "Outputs information about milpa and known repos. See milpa help itself doctor",
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		RunE: func(_ *cobra.Command, args []string) (err error) {
			bold := color.New(color.Bold)
			warn := color.New(color.FgYellow)
			fail := color.New(color.FgRed)
			success := color.New(color.FgGreen)
			failedOverall := false

			var milpaRoot string
			if mp := os.Getenv("MILPA_ROOT"); mp != "" {
				milpaRoot = strings.Join(strings.Split(mp, ":"), "\n")
			} else {
				milpaRoot = warn.Sprint("empty")
			}
			bold.Printf("MILPA_ROOT is: %s\n", milpaRoot)

			var milpaPath string
			bold.Printf("MILPA_PATH is: ")
			if mp := os.Getenv("MILPA_PATH"); mp != "" {
				milpaPath = "\n" + strings.Join(strings.Split(mp, ":"), "\n")
			} else {
				milpaPath = warn.Sprint("empty")
			}
			fmt.Printf("%s\n", milpaPath)
			fmt.Println("")
			bold.Printf("Runnable commands:\n")

			sort.Sort(ByPath(commands))
			for _, cmd := range commands {
				report := cmd.Validate()
				message := ""

				hasFailures := false
				for property, isValid := range report {
					formatter := success
					if !isValid {
						hasFailures = true
						formatter = fail
					}

					message += formatter.Sprintf("  - %s\n", property)
				}
				prefix := "✅"
				if hasFailures {
					failedOverall = true
					prefix = "❌"
				}

				fmt.Println(bold.Sprintf("%s %s", prefix, cmd.FullName()), "—", cmd.Meta.Path)
				if message != "" {
					fmt.Println(message)
				}
				fmt.Println("-----------")
			}

			if failedOverall {
				return fmt.Errorf("your milpa could use some help, check out errors above")
			}

			return
		},
	}
}

func writeDocs(dst string, path []string, cmd *cobra.Command) error {
	if !cmd.IsAvailableCommand() {
		return nil
	}

	dir := strings.Join(append([]string{dst}, path...), "/")
	name := cmd.Name()

	if cmd.HasAvailableSubCommands() {
		dir = dir + "/" + name
		name = "_index"
	}

	logrus.Debugf("Creating directory %s", dir)
	os.MkdirAll(dir, 0760)
	fname := dir + "/" + name + ".md"

	frontMatter := `---
title: ` + strings.Join(append(path, cmd.Name()), " ") + `
type: docs
---
`

	var tmp bytes.Buffer
	cmd.SetOutput(&tmp)
	cmd.Help()

	logrus.Debugf("Creating file %s", fname)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	fixedLinks := bytes.ReplaceAll(tmp.Bytes(), []byte("[.milpa/"), []byte("[/"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("index.md"), []byte(""))

	_, err = f.Write(append([]byte(frontMatter), fixedLinks...))
	if err != nil {
		return err
	}

	if cmd.HasAvailableSubCommands() {
		for _, cc := range cmd.Commands() {
			err := writeDocs(dst, append(path, cmd.Name()), cc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var generateDocumentationCommand *cobra.Command = &cobra.Command{
	Use:               "__generate_documentation [DST]",
	Short:             "Outputs markdownd documentation for all known commands",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		path := []string{}
		dst := args[0]

		err = writeDocs(dst, path, cmd.Root())
		if err != nil {
			return err
		}

		return nil
	},
}
