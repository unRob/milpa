// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"git.rob.mx/nidito/chinampa/pkg/tree"
	milpaCmd "github.com/unrob/milpa/internal/command"

	"gopkg.in/yaml.v3"
)

var ctLog = logger.Sub("itself command-tree")

type serializar func(interface{}) ([]byte, error)

func addMetaToTree(t *tree.CommandTree) {
	if t.Command != nil && t.Command.Meta == nil {
		meta := &milpaCmd.Meta{
			Path: "",
			Repo: "",
			Name: t.Command.Path,
			Kind: milpaCmd.KindVirtual,
		}
		t.Command.Meta = &meta
	}

	for _, subT := range t.Children {
		addMetaToTree(subT)
	}
}

var CommandTree = &command.Command{
	Path:    []string{"__command_tree"},
	Hidden:  true,
	Summary: "Outputs a representation of the command tree at a given PREFIX",
	Description: `Prints out command names and descriptions, or optionally all properties as ﹅json﹅ or ﹅yaml﹅.

  ## Examples

  ﹅﹅﹅sh
  # print all known subcommands
  ` + runtime.Executable + ` __command_tree

  # print a tree of "` + runtime.Executable + ` itself" sub-commands
  ` + runtime.Executable + ` __command_tree

  # print out all commands, skipping groups
  ` + runtime.Executable + ` __command_tree --template '{{ if (not (eq .Meta.Kind "virtual")) }}{{ .FullName }}'$'\n''{{ end }}'

  # get all commands as a json tree
  ` + runtime.Executable + ` __command_tree --output json

  # same, but as the yaml representation of this command itself
  ` + runtime.Executable + ` __command_tree --output json itself command-tree
  ﹅﹅﹅`,
	Arguments: command.Arguments{
		{
			Name:        "prefix",
			Description: "The prefix to traverse the command tree from",
			Variadic:    true,
			Default:     []string{},
		},
	},
	Options: command.Options{
		"depth": &command.Option{
			Default:     15,
			Type:        command.ValueTypeInt,
			Description: "The maximum depth to search for commands",
		},
		"format": &command.Option{
			Default:     "text",
			Description: "The format to output results in",
			Values: &command.ValueSource{
				Static: &([]string{"yaml", "json", "text", "autocomplete"}),
			},
		},
		"template": &command.Option{
			Default:     "{{ .Name }} - {{ .Summary }}\n",
			Description: "a go-template to apply to every command",
		},
	},
	Action: func(cmd *command.Command) error {
		args := cmd.Arguments[0].ToValue().([]string)

		base, remainingArgs, err := cmd.Cobra.Root().Find(args)
		if err != nil {
			return err
		}

		if len(remainingArgs) > 0 && len(remainingArgs) == len(args) {
			return nil
		}

		depth := cmd.Options["depth"].ToValue().(int)
		format := cmd.Options["format"].ToString()

		ctLog.Debugf("looking for commands at %s depth: %d", base.Name(), depth)
		tree.Build(base, depth)

		var serializationFn serializar
		addMeta := func(res serializar) serializar {
			return func(i interface{}) ([]byte, error) {
				if t, ok := i.(*tree.CommandTree); ok {
					addMetaToTree(t)
				}
				return res(i)
			}
		}
		switch format {
		case "yaml":
			serializationFn = addMeta(yaml.Marshal)
		case "json":
			serializationFn = addMeta(func(t interface{}) ([]byte, error) { return json.MarshalIndent(t, "", "  ") })
		case "text":
			outputTpl := cmd.Options["template"].ToString()

			tpl := template.Must(template.New("treeItem").Funcs(render.TemplateFuncs).Parse(outputTpl))
			serializationFn = func(t interface{}) ([]byte, error) {
				tree := t.(*tree.CommandTree)
				var output bytes.Buffer
				err := tree.Traverse(func(cmd *command.Command) error { return tpl.Execute(&output, cmd) })
				return output.Bytes(), err
			}
		case "autocomplete":
			serializationFn = func(interface{}) ([]byte, error) {
				return []byte(strings.Join(tree.ChildrenNames(), "\n") + "\n"), nil
			}
		default:
			return errors.BadArguments{Msg: fmt.Sprintf("Unknown format <%s> for command tree serialization", format)}
		}

		serialized, err := tree.Serialize(serializationFn)
		if err != nil {
			return err
		}
		fmt.Print(serialized)

		return nil
	},
}
