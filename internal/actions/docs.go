// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/statuscode"
	"github.com/spf13/cobra"
	"github.com/unrob/milpa/internal/bootstrap"

	milpaCommand "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/logger"
	"github.com/unrob/milpa/internal/lookup"
)

var dlog = logger.Sub("action:docs")

func readDoc(query []string) ([]byte, error) {
	if err := bootstrap.CheckMilpaPathSet(); err != nil {
		return []byte{}, err
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("requesting docs help")
	}

	queryString := strings.Join(query, "/")

	for _, path := range bootstrap.MilpaPath {
		candidate := path + "/docs/" + queryString
		dlog.Debugf("looking for doc named %s", candidate)
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

	dlog.Debugf("Creating directory %s", dir)
	if err := os.MkdirAll(dir, 0760); err != nil {
		return err
	}
	fname := dir + "/" + name + ".md"

	var tmp bytes.Buffer
	cmd.SetOutput(&tmp)

	if err := cmd.Help(); err != nil {
		return err
	}

	dlog.Debugf("Creating file %s", fname)
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

var Docs = &command.Command{
	Path:        []string{"help", "docs"},
	Summary:     "Dislplays docs on TOPIC",
	Description: "Shows markdown-formatted documentation from milpa repos. See `" + _c.Milpa + " " + _c.HelpCommandName + " docs milpa repo docs` for more information on how to write your own.",
	Arguments: command.Arguments{
		&command.Argument{
			Name:        "topic",
			Description: "The topic to show docs for",
			Variadic:    true,
			Required:    true,
			Values: &command.ValueSource{
				Suggestion: true,
				Func: func(cmd *command.Command, currentValue, config string) (values []string, flag cobra.ShellCompDirective, err error) {
					args := cmd.Arguments[0].ToValue().([]string)
					dlog.Debugf("looking for docs given %v and %s", args, currentValue)

					cv := ""
					if len(args) > 1 {
						cv = args[len(args)-1]
						args = args[0 : len(args)-1]
					}
					dlog.Debugf("looking for docs given %v and %s", args, cv)
					docs, err := lookup.Docs(args, cv, false)
					if err != nil {
						return nil, cobra.ShellCompDirectiveNoFileComp, err
					}

					return docs, cobra.ShellCompDirectiveNoFileComp, nil
				},
			},
		},
	},
	Meta: milpaCommand.Meta{
		Path: os.Getenv(_c.EnvVarMilpaRoot) + "/milpa/docs",
		Name: []string{_c.HelpCommandName, "docs"},
		Repo: os.Getenv(_c.EnvVarMilpaRoot),
		Kind: "docs",
	},
	HelpFunc: func(printLinks bool) string {
		dlog.Debug("showing docs help")
		topics, err := lookup.Docs([]string{}, "", false)
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
	Action: func(cmd *command.Command) error {
		dlog.Debug("Rendering docs")
		args := cmd.Arguments[0].ToValue().([]string)
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

		withColor, _ := cmd.Cobra.Flags().GetBool("no-color")

		doc, err := render.Markdown(contents, !withColor)
		if err != nil {
			return err
		}

		if _, err := cmd.Cobra.OutOrStderr().Write(doc); err != nil {
			return err
		}
		os.Exit(statuscode.RenderHelp)

		return nil
	},
}

func fixLinks(contents []byte) []byte {
	fixedLinks := bytes.ReplaceAll(contents, []byte("(/"+_c.RepoDocs), []byte("(/help/docs"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("(/"+_c.RepoCommands+"/"), []byte("(/"))
	fixedLinks = bytes.ReplaceAll(fixedLinks, []byte("index.md"), []byte(""))
	return bytes.ReplaceAll(fixedLinks, []byte(".md"), []byte("/"))
}

func writeDocs(dst string) error {
	allDocs, err := lookup.AllDocs()
	if err != nil {
		return err
	}

	for _, doc := range allDocs {
		contents, err := os.ReadFile(doc)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(doc, ".md")
		for _, p := range bootstrap.MilpaPath {
			name = strings.Replace(name, p, "", 1)
		}

		components := strings.Split(name, "/")
		last := len(components) - 1
		dir := fmt.Sprintf("%s/help%s", dst, strings.Join(components[0:last], "/"))

		if components[last] == "index" {
			components[last] = "_index"
		}

		fname := fmt.Sprintf("%s/help%s.md", dst, strings.Join(components, "/"))
		dlog.Debugf("Creating dir %s", dir)
		if err := os.MkdirAll(dir, 0760); err != nil {
			return err
		}

		dlog.Debugf("Creating file %s (%s)", fname, dir)
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

var GenerateDocs = &command.Command{
	Path:        []string{"__generate_documentation"},
	Hidden:      true,
	Summary:     "Outputs markdown documentation for all known commands",
	Description: "Creates a set of nested folders at `DEST` with markdown files for every command.",
	Action: func(cmd *command.Command) error {
		path := []string{}
		dst := cmd.Arguments[0].ToString()

		return writeCommandDocs(dst, path, cmd.Cobra.Root())
	},
}
