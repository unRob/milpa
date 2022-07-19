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
package actions

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	command "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/registry"
	"github.com/unrob/milpa/internal/render"
	runtime "github.com/unrob/milpa/internal/runtime"
)

func readDoc(query []string) ([]byte, error) {
	if err := runtime.CheckMilpaPathSet(); err != nil {
		return []byte{}, err
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("requesting docs help")
	}

	queryString := strings.Join(query, "/")

	for _, path := range runtime.MilpaPath {
		candidate := path + "/docs/" + queryString
		logrus.Debugf("looking for doc named %s", candidate)
		_, err := os.Lstat(candidate + ".md")
		if err == nil {
			return os.ReadFile(candidate + ".md")
		}

		if _, err := os.Lstat(candidate + "/index.md"); err == nil {
			return os.ReadFile(candidate + "/index.md")
		}

		if _, err := os.Stat(candidate); err == nil {
			return []byte{}, errors.BadArguments{Msg: fmt.Sprintf("Missing topic for %s", strings.Join(query, " "))}
		}
	}

	return nil, fmt.Errorf("doc not found")
}

func writeCommandDocs(dst string, path []string, cmd *cobra.Command) error {
	if !cmd.IsAvailableCommand() && cmd.Name() != _c.HelpCommandName {
		return nil
	}

	dir := strings.Join(append([]string{dst}, path...), "/")
	name := cmd.Name()

	if cmd.HasAvailableSubCommands() {
		if name != _c.Milpa {
			dir = dir + "/" + name
		}
		name = "_index"
	}

	logrus.Debugf("Creating directory %s", dir)
	if err := os.MkdirAll(dir, 0760); err != nil {
		return err
	}
	fname := dir + "/" + name + ".md"

	var tmp bytes.Buffer
	cmd.SetOutput(&tmp)

	if err := cmd.Help(); err != nil {
		return err
	}

	logrus.Debugf("Creating file %s", fname)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(fixLinks(tmp.Bytes()))
	if err != nil {
		return err
	}

	if cmd.HasAvailableSubCommands() {
		for _, cc := range cmd.Commands() {
			subPath := path
			if cmd.Name() != "milpa" {
				subPath = append(path, cmd.Name()) //nolint:gocritic
			}

			err := writeCommandDocs(dst, subPath, cc)
			if err != nil {
				return err
			}
		}
	} else if cmd.Annotations["MilpaDocs"] == "true" {
		err := writeDocs(dst)
		if err != nil {
			return err
		}
	}

	return nil
}

var docs = &command.Command{
	Summary:     "Dislplays docs on TOPIC",
	Description: "Shows markdown-formatted documentation from milpa repos. See `" + _c.Milpa + " " + _c.HelpCommandName + " docs milpa repo docs` for more information on how to write your own.",
	Arguments: command.Arguments{
		&command.Argument{
			Name:        "topic",
			Description: "The topic to show docs for",
			Variadic:    true,
			Required:    true,
		},
	},
	Meta: command.Meta{
		Path: os.Getenv(_c.EnvVarMilpaRoot) + "/milpa/docs",
		Name: []string{_c.HelpCommandName, "docs"},
		Repo: os.Getenv(_c.EnvVarMilpaRoot),
		Kind: "internal",
	},
	HelpFunc: func(printLinks bool) string {
		topics, err := registry.FindDocs([]string{}, "", false)
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

func fixLinks(contents []byte) []byte {
	fixedLinks := bytes.ReplaceAll(contents, []byte("(/"+_c.RepoDocs), []byte("(/help/docs"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("(/"+_c.RepoCommands+"/"), []byte("(/"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("index.md"), []byte(""))
	return bytes.ReplaceAll(fixedLinks, []byte(".md"), []byte("/"))
}

func writeDocs(dst string) error {
	allDocs, err := registry.FindAllDocs()
	if err != nil {
		return err
	}

	for _, doc := range allDocs {
		contents, err := os.ReadFile(doc)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(doc, ".md")
		for _, p := range runtime.MilpaPath {
			name = strings.Replace(name, p, "", 1)
		}

		components := strings.Split(name, "/")
		last := len(components) - 1
		dir := fmt.Sprintf("%s/help%s", dst, strings.Join(components[0:last], "/"))

		if components[last] == "index" {
			components[last] = "_index"
		}

		fname := fmt.Sprintf("%s/help%s.md", dst, strings.Join(components, "/"))
		logrus.Debugf("Creating dir %s", dir)
		if err := os.MkdirAll(dir, 0760); err != nil {
			return err
		}

		logrus.Debugf("Creating file %s (%s)", fname, dir)
		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(fixLinks(contents))
		if err != nil {
			return err
		}
	}

	return nil
}

var docsCommand = &cobra.Command{
	Use:   "docs [TOPIC]",
	Short: docs.Summary,
	Long:  docs.Description,
	ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		logrus.Debugf("looking for docs given %v and %s", args, toComplete)
		docs, err := registry.FindDocs(args, toComplete, false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return docs, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.BadArguments{Msg: "Missing doc topic to display"}
		}

		contents, err := readDoc(args)
		if err != nil {
			switch err.(type) {
			case errors.BadArguments:
				return err
			}
			return errors.NotFound{Msg: "Unknown doc: " + err.Error()}
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

		withColor, _ := c.Flags().GetBool("no-color")

		doc, err := render.Markdown(contents, !withColor)
		if err != nil {
			return err
		}

		if _, err := c.OutOrStderr().Write(doc); err != nil {
			return err
		}
		os.Exit(_c.ExitStatusRenderHelp)

		return nil
	},
	Annotations: map[string]string{
		"MilpaDocs": "true",
	},
}

var generateDocumentationCommand = &cobra.Command{
	Use:               "__generate_documentation [DST]",
	Short:             "Outputs markdown documentation for all known commands",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		path := []string{}
		dst := args[0]

		err = writeCommandDocs(dst, path, cmd.Root())
		if err != nil {
			return err
		}

		return nil
	},
}
