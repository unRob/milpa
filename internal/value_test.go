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
package internal_test

import (
	"testing"

	"github.com/spf13/pflag"
	. "github.com/unrob/milpa/internal"
)

func TestResolveTemplate(t *testing.T) {
	cmd := &Command{
		Arguments: []Argument{
			{
				Name:    "argument_0",
				Default: "default",
			},
			{
				Name:     "argument_n",
				Variadic: true,
			},
		},
		Options: Options{
			"option": {
				Default: "default",
				Type:    "string",
			},
			"bool": {
				Type:    "bool",
				Default: false,
			},
		},
	}

	overrideFlags := &pflag.FlagSet{}
	overrideFlags.String("option", "override", "stuff")
	overrideFlags.Bool("bool", false, "stuff")
	overrideFlags.Bool("help", false, "stuff")
	overrideFlags.Bool("no-color", false, "stuff")
	err := overrideFlags.Parse([]string{"--option", "override", "--bool", "--help", "--no-color"})
	if err != nil {
		t.Fatalf("Could not parse test flags")
	}

	cases := []struct {
		Tpl      string
		Expected string
		Args     []string
		Flags    *pflag.FlagSet
		Errors   bool
	}{
		{
			Tpl:      "adds nothing to nothing",
			Expected: "adds nothing to nothing",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `prints default option as {{ Opt "option" }}`,
			Expected: "prints default option as --option default",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `prints default option value as {{ .Opts.option }}`,
			Expected: "prints default option value as default",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `prints default argument as {{ Arg "argument_0" }}`,
			Expected: "prints default argument as default",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `prints default argument value as {{ .Args.argument_0 }}`,
			Expected: "prints default argument value as default",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `overrides default option as {{ Opt "option" }}`,
			Expected: "overrides default option as --option override",
			Errors:   false,
			Args:     []string{},
			Flags:    overrideFlags,
		},
		{
			Tpl:      `overrides default argument as {{ Arg "argument_0" }}`,
			Expected: "overrides default argument as override",
			Errors:   false,
			Args:     []string{"override"},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `combines defaults as {{ Opt "option" }} {{ Opt "bool"}} {{ Arg "argument_0" }}`,
			Expected: "combines defaults as --option default --bool false default",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `combines overrides as {{ Opt "option" }} {{ Opt "bool" }} {{ Arg "argument_0" }}`,
			Expected: "combines overrides as --option override --bool true twice",
			Errors:   false,
			Args:     []string{"twice"},
			Flags:    overrideFlags,
		},
		{
			Tpl:      `prints variadic as {{ Arg "argument_0" }} {{ Arg "argument_n" }}`,
			Expected: "prints variadic as override a b",
			Errors:   false,
			Args:     []string{"override", "a", "b"},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:      `doesn't error on bad names {{ Opt "bad-option" }} {{ Arg "bad-argument" }}`,
			Expected: "doesn't error on bad names  ",
			Errors:   false,
			Args:     []string{},
			Flags:    &pflag.FlagSet{},
		},
		{
			Tpl:    `errors on bad templates {{ BadFunc }}`,
			Args:   []string{},
			Flags:  &pflag.FlagSet{},
			Errors: true,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.Expected, func(t *testing.T) {
			res, err := ResolveTemplate(cmd, test.Tpl, test.Args, test.Flags)

			if err != nil && !test.Errors {
				t.Fatalf("good template failed: %s", err)
			}

			if res != test.Expected {
				t.Fatalf("expected '%s' got '%s'", test.Expected, res)
			}
		})
	}

}
