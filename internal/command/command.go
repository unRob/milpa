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
package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	runtime "github.com/unrob/milpa/internal/runtime"
	"gopkg.in/yaml.v2"
)

type Kind string

const (
	KindUnknown    Kind = ""
	KindExecutable Kind = "executable"
	KindSource     Kind = "source"
	KindVirtual    Kind = "virtual"
	KindRoot       Kind = "root"
)

type Command struct {
	Meta         Meta      `json:"meta" yaml:"meta"`
	Summary      string    `json:"summary" yaml:"summary" validate:"required"`
	Description  string    `json:"description" yaml:"description" validate:"required"`
	Arguments    Arguments `json:"arguments" yaml:"arguments" validate:"dive"`
	Options      Options   `json:"options" yaml:"options" validate:"dive"`
	runtimeFlags *pflag.FlagSet
	issues       []string
	HelpFunc     func(printLinks bool) string `json:"-" yaml:"-"`
	cc           *cobra.Command
}

type Meta struct {
	Path string   `json:"path" yaml:"path"`
	Repo string   `json:"repo" yaml:"repo"`
	Name []string `json:"name" yaml:"name"`
	Kind Kind     `json:"kind" yaml:"kind"`
}

func metaForPath(path string, repo string) (meta Meta) {
	meta.Path = path
	meta.Repo = repo
	name := strings.TrimSuffix(path, ".sh")
	name = strings.TrimPrefix(name, repo+"/"+_c.RepoCommandFolderName+"/")
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
		err = errors.ConfigError{
			Err:    err,
			Config: spec,
		}
		cmd.issues = append(cmd.issues, err.Error())
	}

	return cmd.SetBindings(), nil
}

func (cmd *Command) SetBindings() *Command {
	ptr := cmd
	for _, opt := range cmd.Options {
		opt.Command = ptr
		if opt.Validates() {
			opt.Values.Command = ptr
		}
	}

	for _, arg := range cmd.Arguments {
		arg.Command = ptr
		if arg.Validates() {
			arg.Values.Command = ptr
		}
	}
	return ptr
}

func (cmd *Command) Name() string {
	return cmd.Meta.Name[len(cmd.Meta.Name)-1]
}

func (cmd *Command) FullName() string {
	return strings.Join(cmd.Meta.Name, " ")
}

func (cmd *Command) FlagSet() *pflag.FlagSet {
	if cmd.runtimeFlags == nil {
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
	return cmd.runtimeFlags
}

func (cmd *Command) Run(cc *cobra.Command, args []string) error {
	logrus.Debugf("running command %s", cmd.FullName())
	cmd.Arguments.Parse(args)
	skipValidation, _ := cc.Flags().GetBool("skip-validation")
	cmd.Options.Parse(cc.Flags())
	if !skipValidation && runtime.ValidationEnabled() {
		logrus.Debug("Validating arguments")
		if err := cmd.Arguments.AreValid(); err != nil {
			return err
		}

		logrus.Debug("Validating flags")
		if err := cmd.Options.AreValid(); err != nil {
			return err
		}
	}

	env := cmd.ToEval(args)

	if os.Getenv(_c.EnvVarCompaOut) != "" {
		return os.WriteFile(os.Getenv(_c.EnvVarCompaOut), []byte(env), 0600)
	}

	fmt.Println(env)
	return nil
}

func (cmd *Command) RunStandAlone(cc *cobra.Command, args []string) error {
	cmd.Arguments.Parse(args)
	skipValidation, _ := cc.Flags().GetBool("skip-validation")
	cmd.Options.Parse(cc.Flags())
	if !skipValidation && runtime.ValidationEnabled() {
		if err := cmd.Arguments.AreValid(); err != nil {
			return err
		}
		if err := cmd.Options.AreValid(); err != nil {
			return err
		}
	}

	env := cmd.Env(os.Environ())

	entrypoint := ""
	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}

	if cmd.Meta.Kind == KindSource {
		entrypoint = fmt.Sprintf("%s/%s", runtime.MilpaRoot, "milpa")
	}

	child := exec.Command(entrypoint, args...)
	child.Env = env
	child.Stderr = os.Stderr
	child.Stdout = os.Stdout
	child.Stdin = os.Stdin

	logrus.Debugf("Found command of kind %s at %s", cmd.Meta.Kind, entrypoint)

	if err := child.Run(); err != nil {
		exitCode := child.ProcessState.ExitCode()
		return errors.SubCommandExit{Err: err, ExitCode: exitCode}
	}

	return nil
}

func (cmd *Command) SetCobra(cc *cobra.Command) {
	cmd.cc = cc
}
