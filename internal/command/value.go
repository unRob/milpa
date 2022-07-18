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
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/exec"
)

// ValueType represent the kinds of or option.
type ValueType string

const (
	// ValueTypeDefault is the empty string, maps to ValueTypeString.
	ValueTypeDefault ValueType = ""
	// ValueTypeString a value treated like a string.
	ValueTypeString ValueType = "string"
	// ValueTypeBoolean is a value treated like a boolean.
	ValueTypeBoolean ValueType = "bool"
)

// ValueSource represents the source for an auto-completed and/or validated option/argument.
type ValueSource struct {
	// Directories prompts for directories with the given prefix.
	Directories *string `json:"dirs,omitempty" yaml:"dirs,omitempty" validate:"omitempty,excluded_with=Files Script Static Milpa"`
	// Files prompts for files with the given extensions
	Files *[]string `json:"files,omitempty" yaml:"files,omitempty" validate:"omitempty,excluded_with=Directories Script Static Milpa"`
	// Script runs the provided command with `bash -c "$script"` and returns an option for every line of stdout.
	Script string `json:"script,omitempty" yaml:"script,omitempty" validate:"omitempty,excluded_with=Directories Files Static Milpa"`
	// Static returns the given list.
	Static *[]string `json:"static,omitempty" yaml:"static,omitempty" validate:"omitempty,excluded_with=Directories Files Script Milpa"`
	// Milpa runs a subcommand and returns an option for every line of stdout.
	Milpa string `json:"milpa,omitempty" yaml:"milpa,omitempty" validate:"omitempty,excluded_with=Directories Files Script Static"`
	// Timeout is the maximum amount of time milpa will wait for a Script or Milpa command before giving up on completions/validations.
	Timeout int `json:"timeout,omitempty" yaml:"timeout,omitempty" validate:"omitempty,excluded_with=Directories Files Static"`
	// Suggestion if provided will only suggest autocomplete values but will not perform validation of a given value
	Suggestion bool     `json:"suggest-only" yaml:"suggest-only" validate:"omitempty"` // nolint:tagliatelle
	Command    *Command `json:"-" yaml:"-" validate:"-"`
	computed   *[]string
	flag       cobra.ShellCompDirective
}

// Validates tells if a value needs to be validated.
func (vs *ValueSource) Validates() bool {
	if vs.Directories != nil || vs.Files != nil {
		return false
	}

	return !vs.Suggestion
}

// Resolve returns the values for autocomplete and validation.
func (vs *ValueSource) Resolve() (values []string, flag cobra.ShellCompDirective, err error) {
	if vs.computed != nil {
		return *vs.computed, vs.flag, nil
	}

	if vs.Timeout == 0 {
		vs.Timeout = 5
	}

	flag = cobra.ShellCompDirectiveDefault
	timeout := time.Duration(vs.Timeout)

	switch {
	case vs.Static != nil:
		values = *vs.Static
	case vs.Files != nil:
		flag = cobra.ShellCompDirectiveFilterFileExt
		values = *vs.Files
	case vs.Directories != nil:
		flag = cobra.ShellCompDirectiveFilterDirs
		values = []string{*vs.Directories}
	case vs.Milpa != "" || vs.Script != "":
		if vs.Command == nil {
			return nil, cobra.ShellCompDirectiveError, fmt.Errorf("bug: command is nil")
		}
		cmd, err := vs.Command.ResolveTemplate(vs.Milpa + vs.Script)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError, err
		}

		var args []string
		if vs.Script != "" {
			args = append([]string{"/bin/bash", "-c"}, cmd)
		} else {
			args = append([]string{"milpa"}, strings.Split(cmd, " ")...)
		}

		envMap := vs.Command.EnvironmentMap()
		env := os.Environ()
		for k, v := range envMap {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		values, flag, err = exec.Exec(vs.Command.FullName(), args, env, timeout*time.Second)
		if err != nil {
			return nil, flag, err
		}
	}

	vs.computed = &values
	vs.flag = flag

	return values, flag, err
}

type AutocompleteTemplate struct {
	Args map[string]string
	Opts map[string]string
}

func (tpl *AutocompleteTemplate) Opt(name string) string {
	if val, ok := tpl.Opts[name]; ok {
		return fmt.Sprintf("--%s %s", name, val)
	}

	return ""
}

func (tpl *AutocompleteTemplate) Arg(name string) string {
	return tpl.Args[name]
}

func (cmd *Command) ResolveTemplate(templateString string) (string, error) {
	var buf bytes.Buffer

	tplData := &AutocompleteTemplate{
		Args: cmd.Arguments.AllKnown(),
		Opts: cmd.Options.AllKnown(),
	}

	fnMap := template.FuncMap{
		"Opt": tplData.Opt,
		"Arg": tplData.Arg,
	}

	for k, v := range _c.TemplateFuncs {
		fnMap[k] = v
	}

	tpl, err := template.New("subcommand").Funcs(fnMap).Parse(templateString)

	if err != nil {
		return "", err
	}

	err = tpl.Execute(&buf, tplData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
