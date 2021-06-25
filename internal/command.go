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
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Command struct {
	Meta         CommandMeta              `json:"_meta"`
	Summary      string                   `yaml:"summary"`
	Description  string                   `yaml:"description"`
	Arguments    CommandArguments         `yaml:"arguments"`
	Options      map[string]CommandOption `yaml:"options"`
	runtimeFlags *pflag.FlagSet
}

type CommandMeta struct {
	Path    string
	Package string
	Name    []string
	Kind    string
}

func New(path string, spec string, pkg string, kind string) (*Command, error) {
	contents, err := ioutil.ReadFile(spec)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", spec, err)
	}

	commandPath := strings.SplitN(path, fmt.Sprintf("%s/.milpa/commands/", pkg), 2)[1]
	commandName := strings.Split(strings.TrimSuffix(commandPath, ".sh"), "/")

	cmd, err := parseCommand(contents, CommandMeta{
		Path:    path,
		Package: pkg,
		Name:    commandName,
		Kind:    kind,
	})
	if err != nil {
		err = fmt.Errorf("error parsing %s: %w", spec, err)
	}
	return cmd, err
}

func parseCommand(yamlBytes []byte, meta CommandMeta) (*Command, error) {
	cmd := &Command{
		Meta:      meta,
		Arguments: []CommandArgument{},
		Options:   map[string]CommandOption{},
	}
	err := yaml.Unmarshal(yamlBytes, cmd)

	if err != nil {
		if err, ok := err.(ConfigError); ok {
			err.Config = meta.Path
		}
		return nil, err
	}

	return cmd, nil
}

func (cmd *Command) FullName() string {
	return strings.Join(cmd.Meta.Name, " ")
}

func (cmd *Command) CreateFlagSet() error {
	if cmd.runtimeFlags != nil {
		return nil
	}
	fs := pflag.NewFlagSet(strings.Join(cmd.Meta.Name, " "), pflag.ContinueOnError)
	fs.SortFlags = false
	fs.Usage = func() {}

	for name, opt := range cmd.Options {
		switch opt.Type {
		case ValueTypeBoolean:
			def := false
			if opt.Default != nil {
				def = opt.Default.(bool)
			}
			fs.Bool(name, def, opt.Description)
		case ValueTypeDefault, ValueTypeString:
			opt.Type = ValueTypeString
			def := ""
			if opt.Default != nil {
				def = opt.Default.(string)
			}
			fs.String(name, def, opt.Description)
		default:
			return fmt.Errorf("unknown option type: <%s> for option: <%s>", opt.Type, name)
		}
	}

	cmd.runtimeFlags = fs
	return nil
}
