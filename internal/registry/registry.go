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
package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	command "github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
)

var registry = &CommandRegistry{
	kv: map[string]*command.Command{},
}

type ByPath []*command.Command

func (cmds ByPath) Len() int      { return len(cmds) }
func (cmds ByPath) Swap(i, j int) { cmds[i], cmds[j] = cmds[j], cmds[i] }
func (cmds ByPath) Less(i, j int) bool {
	if cmds[i].Meta.Path == cmds[j].Meta.Path {
		return cmds[i].FullName() < cmds[j].FullName()
	}
	return cmds[i].Meta.Path < cmds[j].Meta.Path
}

type CommandTree struct {
	Command  *command.Command `json:"command"`
	Children []*CommandTree   `json:"children"`
}

type CommandRegistry struct {
	kv     map[string]*command.Command
	byPath []*command.Command
	tree   *CommandTree
}

func Register(cmd *command.Command) {
	registry.kv[cmd.FullName()] = cmd
}

func Get(id string) *command.Command {
	return registry.kv[id]
}

func CommandList() []*command.Command {
	if len(registry.byPath) == 0 {
		list := []*command.Command{}
		for _, v := range registry.kv {
			list = append(list, v)
		}
		sort.Sort(ByPath(list))
		registry.byPath = list
	}

	return registry.byPath
}

func BuildTree(cc *cobra.Command, depth int) {
	tree := &CommandTree{
		Command:  fromCobra(cc),
		Children: []*CommandTree{},
	}

	var populateTree func(cmd *cobra.Command, ct *CommandTree, maxDepth int, depth int)
	populateTree = func(cmd *cobra.Command, ct *CommandTree, maxDepth int, depth int) {
		newDepth := depth + 1
		for _, subcc := range cmd.Commands() {
			if subcc.Hidden {
				continue
			}

			if cmd := fromCobra(subcc); cmd != nil {
				leaf := &CommandTree{Children: []*CommandTree{}}
				leaf.Command = cmd
				ct.Children = append(ct.Children, leaf)

				if newDepth < maxDepth {
					populateTree(subcc, leaf, maxDepth, newDepth)
				}
			}
		}
	}
	populateTree(cc, tree, depth, 0)

	registry.tree = tree
}

func AsJSONTree() (string, error) {
	bytes, err := json.Marshal(registry.tree)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ChildrenNames() []string {
	if registry.tree == nil {
		return []string{}
	}

	ret := make([]string, len(registry.tree.Children))
	for idx, cmd := range registry.tree.Children {
		ret[idx] = cmd.Command.Meta.Name[len(cmd.Command.Meta.Name)-1]
	}
	return ret
}

func SetRoot(ccRoot *cobra.Command, cmdRoot *command.Command) {
	for _, cmd := range CommandList() {
		cmd := cmd
		leaf := toCobra(cmd, cmdRoot.Options)
		container := ccRoot
		for idx, cp := range cmd.Meta.Name {
			if idx == len(cmd.Meta.Name)-1 {
				logrus.Debugf("adding command %s to %s", leaf.Name(), container.Name())
				container.AddCommand(leaf)
				break
			}

			query := []string{cp}
			if cc, _, err := container.Find(query); err == nil && cc != container {
				logrus.Debugf("found %s in %s", query, cc.Name())
				container = cc
			} else {
				logrus.Debugf("creating %s in %s", query, container.Name())
				groupName := strings.Join(query, " ")
				cc := &cobra.Command{
					Use:                        cp,
					Short:                      fmt.Sprintf("%s subcommands", groupName),
					DisableAutoGenTag:          true,
					SuggestionsMinimumDistance: 2,
					SilenceUsage:               true,
					SilenceErrors:              true,
					Annotations: map[string]string{
						_c.ContextKeyRuntimeIndex: groupName,
					},
					Args: func(cmd *cobra.Command, args []string) error {
						err := cobra.OnlyValidArgs(cmd, args)
						if err != nil {

							suggestions := []string{}
							bold := color.New(color.Bold)
							for _, l := range cmd.SuggestionsFor(args[len(args)-1]) {
								suggestions = append(suggestions, bold.Sprint(l))
							}
							last := len(args) - 1
							parent := cmd.CommandPath()
							errMessage := fmt.Sprintf("Unknown subcommand %s of known command %s", bold.Sprint(args[last]), bold.Sprint(parent))
							if len(suggestions) > 0 {
								errMessage += ". Perhaps you meant " + strings.Join(suggestions, ", ") + "?"
							}
							return errors.NotFound{Msg: errMessage, Group: []string{}}
						}
						return nil
					},
					ValidArgs: []string{""},
					RunE: func(cc *cobra.Command, args []string) error {
						if len(args) == 0 {
							return errors.NotFound{Msg: "No subcommand provided", Group: []string{}}
						}
						os.Exit(_c.ExitStatusNotFound)
						return nil
					},
				}

				pathComps := strings.Split(cmd.Meta.Path, "/")
				groupParent := &command.Command{
					Summary:     fmt.Sprintf("%s subcommands", groupName),
					Description: fmt.Sprintf("Runs subcommands within %s", groupName),
					Arguments:   command.Arguments{},
					Options:     command.Options{},
					Meta: command.Meta{
						Name: query,
						Path: strings.Join(pathComps[0:len(pathComps)-1], "/"),
						Repo: cmd.Meta.Repo,
						Kind: "virtual",
					},
				}
				Register(groupParent)
				cc.SetHelpFunc(groupParent.HelpRenderer(cmd.Options))
				container.AddCommand(cc)
				container = cc
			}
		}
	}
	cmdRoot.SetCobra(ccRoot)
}
