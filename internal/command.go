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

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

const cmdPath = ".milpa/commands"

type Command struct {
	Meta         Meta
	Summary      string    `yaml:"summary" validate:"required"`
	Description  string    `yaml:"description" validate:"required"`
	Arguments    Arguments `yaml:"arguments" validate:"dive"`
	Options      Options   `yaml:"options" validate:"dive"`
	runtimeFlags *pflag.FlagSet
	issues       []string
}

type Meta struct {
	Path string
	Repo string
	Name []string
	Kind string
}

func metaForPath(path string, repo string) (meta Meta) {
	meta.Path = path
	meta.Repo = repo
	name := strings.TrimSuffix(path, ".sh")
	name = strings.TrimPrefix(name, repo+"/"+cmdPath+"/")
	meta.Name = strings.Split(name, "/")

	if strings.HasSuffix(path, ".sh") {
		meta.Kind = "source"
	} else {
		meta.Kind = "exec"
	}

	return
}

func New(path string, repo string, strict bool) (cmd *Command, err error) {
	cmd = &Command{}
	cmd.Meta = metaForPath(path, repo)
	cmd.Arguments = []Argument{}
	cmd.Options = Options{}
	cmd.issues = []string{}

	spec := strings.TrimSuffix(path, ".sh") + ".yaml"
	var contents []byte
	if contents, err = ioutil.ReadFile(spec); err == nil {
		if strict {
			err = yaml.UnmarshalStrict(contents, cmd)
		} else {
			err = yaml.Unmarshal(contents, cmd)
		}
	}

	if err != nil {
		err = ConfigError{
			Err:    err,
			Config: spec,
		}
		cmd.issues = append(cmd.issues, err.Error())
	}

	return
}

func (cmd *Command) Validate() (report map[string]bool) {
	report = map[string]bool{}

	for _, issue := range cmd.issues {
		report[issue] = false
	}

	validate := validator.New()
	err := validate.Struct(cmd)
	if err != nil {
		verrs := err.(validator.ValidationErrors)
		for _, issue := range verrs {
			report[fmt.Sprint(issue)] = false
		}
	}

	return
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
			continue
			// return fmt.Errorf("unknown option type: <%s> for option: <%s>", opt.Type, name)
		}
	}

	cmd.runtimeFlags = fs
	return nil
}