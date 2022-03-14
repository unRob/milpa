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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
	runtime "github.com/unrob/milpa/internal/runtime"
	"gopkg.in/yaml.v2"
)

const cmdPath = "commands"

type Command struct {
	Meta         Meta
	Summary      string    `yaml:"summary" validate:"required"`
	Description  string    `yaml:"description" validate:"required"`
	Arguments    Arguments `yaml:"arguments" validate:"dive"`
	Options      Options   `yaml:"options" validate:"dive"`
	runtimeFlags *pflag.FlagSet
	issues       []string
	helpFunc     func(printLinks bool) string
	cc           *cobra.Command
}

var Root = &Command{
	Summary: "Runs commands found in " + _c.RepoRoot + " folders",
	Description: `﹅milpa﹅ is a command-line tool to care for one's own garden of scripts, its name comes from "milpa", an agricultural method that combines multiple crops in close proximity. You and your team write scripts and a little spec for each command -use bash, or any other language-, and ﹅milpa﹅ provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.

  See [﹅milpa help docs milpa﹅](/.milpa/docs/milpa/index.md) for more information about ﹅milpa﹅`,
	Meta: Meta{
		Path: _c.EnvVarMilpaRoot + "/" + _c.Milpa,
		Name: []string{_c.Milpa},
		Repo: _c.EnvVarMilpaRoot,
		Kind: "root",
	},
	Options: Options{
		_c.HelpCommandName: &Option{
			ShortName:   "h",
			Type:        "bool",
			Description: "Display help for any command",
		},
		"verbose": &Option{
			ShortName:   "v",
			Type:        "bool",
			Default:     runtime.VerboseEnabled(),
			Description: "Log verbose output to stderr",
		},
		"no-color": &Option{
			Type:        "bool",
			Description: "Print to stderr without any formatting codes",
		},
		"silent": &Option{
			Type:        "bool",
			Description: "Silence non-error logging",
		},
		"skip-validation": &Option{
			Type:        "bool",
			Description: "Do not validate any arguments or options",
		},
	},
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
	cmd.Arguments = []*Argument{}
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

func (cmd *Command) FullName() string {
	return strings.Join(cmd.Meta.Name, " ")
}

func (cmd *Command) CreateFlagSet() {
	if cmd.runtimeFlags != nil {
		return
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
				def = fmt.Sprintf("%s", opt.Default)
			}
			fs.String(name, def, opt.Description)
		default:
			// ignore flag
			logrus.Warnf("Ignoring unknown option type <%s> for option <%s>", opt.Type, name)
			continue
		}
	}

	cmd.runtimeFlags = fs
}

type varSearchMap struct {
	Status int
	Name   string
	Usage  string
}

func (cmd *Command) Validate() (report map[string]int) {
	report = map[string]int{}

	for _, issue := range cmd.issues {
		report[issue] = 1
	}

	validate := validator.New()
	err := validate.Struct(cmd)
	if err != nil {
		verrs := err.(validator.ValidationErrors)
		for _, issue := range verrs {
			report[fmt.Sprint(issue)] = 1
		}
	}

	if cmd.Meta.Kind == "source" {
		contents, err := ioutil.ReadFile(cmd.Meta.Path)
		if err != nil {
			report["Could not read source"] = 1
			return
		}

		vars := map[string]map[string]*varSearchMap{
			"argument": {},
			"option":   {},
		}

		for _, arg := range cmd.Arguments {
			vars["argument"][strings.ToUpper(strings.ReplaceAll(arg.Name, "-", "_"))] = &varSearchMap{2, arg.Name, ""}
		}

		for name := range cmd.Options {
			vars["option"][strings.ToUpper(strings.ReplaceAll(name, "-", "_"))] = &varSearchMap{2, name, ""}
		}

		matches := _c.OutputPrefixPattern.FindAllStringSubmatch(string(contents), -1)
		for _, match := range matches {
			varName := match[len(match)-1]
			varKind := match[len(match)-2]

			kind := ""
			if varKind == "OPT" {
				kind = "option"
			} else if varKind == "ARG" {
				kind = "argument"
			}
			haystack := vars[kind]

			_, scriptVarIsValid := haystack[varName]
			if !scriptVarIsValid {
				haystack[varName] = &varSearchMap{Status: 1, Name: varName, Usage: match[0]}
			} else {
				haystack[varName].Status = 0
			}
		}

		for kind, col := range vars {
			for _, thisVar := range col {
				message := ""
				switch thisVar.Status {
				case 0:
					message = fmt.Sprintf("%s '%s' is used", kind, thisVar.Name)
				case 1:
					message = fmt.Sprintf("%s '%s' is used but not defined, declared as '%s'", kind, thisVar.Name, thisVar.Usage)
				case 2:
					message = fmt.Sprintf("%s '%s' is not used but defined", kind, thisVar.Name)
				default:
					message = fmt.Sprintf("Unknown status %d for %s '%s'", thisVar.Status, kind, thisVar.Name)
				}

				report[message] = thisVar.Status
			}
		}
	}

	return report
}
