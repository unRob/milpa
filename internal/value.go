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
	"time"

	"github.com/spf13/cobra"
)

// ValueType represent the kinds of values for an argument or option.
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
	Directories *string `yaml:"dirs" validate:"omitempty,excluded_with=Files Script Static Milpa"`
	// Files prompts for files with the given extensions
	Files *[]string `yaml:"files" validate:"omitempty,excluded_with=Directories Script Static Milpa"`
	// Script runs the provided command with `bash -c "$script"` and returns an option for every line of stdout.
	Script string `yaml:"script" validate:"omitempty,excluded_with=Directories Files Static Milpa"`
	// Static returns the given list.
	Static *[]string `yaml:"static" validate:"omitempty,excluded_with=Directories Files Script Milpa"`
	// Milpa runs a subcommand and returns an option for every line of stdout.
	Milpa string `yaml:"milpa" validate:"omitempty,excluded_with=Directories Files Script Static"`
	// Timeout is the maximum amount of time milpa will wait for a Script or Milpa command before giving up on completions/validations.
	Timeout int `yaml:"timeout" validate:"omitempty,excluded_with=Directories Files Static"`
	// Suggestion if provided will only suggest autocomplete values but will not perform validation of a given value
	Suggestion bool `yaml:"suggest-only" validate:"omitempty"`
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
	case vs.Milpa != "":
		values, flag, err = exec("unknown", vs.Milpa, timeout)
	case vs.Static != nil:
		values = *vs.Static
	case vs.Files != nil:
		flag = cobra.ShellCompDirectiveFilterFileExt
		values = *vs.Files
	case vs.Directories != nil:
		flag = cobra.ShellCompDirectiveFilterDirs
		values = []string{*vs.Directories}
	case vs.Script != "":
		values, flag, err = exec("@bash", vs.Script, timeout)
	}

	if err != nil {
		return
	}
	vs.computed = &values
	vs.flag = flag

	return values, flag, err
}
