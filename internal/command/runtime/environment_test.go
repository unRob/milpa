// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime_test

import (
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/command/runtime"
)

func TestArgumentsToEnv(t *testing.T) {
	cases := []struct {
		Command *command.Command
		Args    []string
		Expect  []string
		Env     []string
	}{
		{
			Args:   []string{"something"},
			Expect: []string{"export MILPA_ARG_FIRST=something"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "required", "present"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:     "first",
						Required: true,
					},
				},
			},
		},
		{
			Args:   []string{},
			Expect: []string{"export MILPA_ARG_FIRST=default"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "default", "present"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
					},
				},
			},
		},
		{
			Args: []string{"zero", "one", "two", "three"},
			Expect: []string{
				"export MILPA_ARG_FIRST=zero",
				"declare -a MILPA_ARG_VARIADIC=( one two three )",
			},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "variadic"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
					},
					{
						Name:     "variadic",
						Variadic: true,
					},
				},
			},
		},
		{
			Args: []string{"zero", "one and stuff", "two", "three"},
			Expect: []string{
				"export MILPA_ARG_FIRST=zero",
				"declare -a MILPA_ARG_VARIADIC=( 'one and stuff' two three )",
			},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "variadic"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
					},
					{
						Name:     "variadic",
						Variadic: true,
					},
				},
			},
		},
		{
			Args: []string{"zero", "one and stuff", "two", "three"},
			Expect: []string{
				"set -x MILPA_ARG_FIRST zero",
				"set -x MILPA_ARG_VARIADIC 'one and stuff' two three",
			},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "variadic"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellFish,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
					},
					{
						Name:     "variadic",
						Variadic: true,
					},
				},
			},
		},
		{
			Args: []string{"zero", "one and stuff", "two", "three"},
			Expect: []string{
				"MILPA_ARG_FIRST=zero",
				"MILPA_ARG_VARIADIC=one and stuff|two|three",
			},
			Command: &command.Command{
				Meta: meta.Meta{
					Name: []string{"test", "variadic"},
					Kind: kind.Executable,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
					},
					{
						Name:     "variadic",
						Variadic: true,
					},
				},
			},
		},
		{
			Args:   []string{},
			Expect: []string{"export MILPA_ARG_FIRST=default"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "static", "default"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
						Values: &command.ValueSource{
							Static: &[]string{
								"default",
								"good",
							},
						},
					},
				},
			},
		},
		{
			Args:   []string{"good"},
			Expect: []string{"export MILPA_ARG_FIRST=good"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "static", "good"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
						Values: &command.ValueSource{
							Static: &[]string{
								"default",
								"good",
							},
						},
					},
				},
			},
		},
		{
			Args:   []string{"good"},
			Expect: []string{"export MILPA_ARG_FIRST=good"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "script", "good"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Arguments: []*command.Argument{
					{
						Name:    "first",
						Default: "default",
						Values: &command.ValueSource{
							Script: "echo good; echo default",
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Command.FullName(), func(t *testing.T) {
			dst := []string{}
			c.Command.SetBindings()
			if err := c.Command.Arguments.Parse(c.Args); err != nil {
				t.Fatal(err)
			}
			runtime.ArgumentsToEnv(c.Command, &dst)

			err := c.Command.Arguments.AreValid()
			if err != nil {
				t.Fatalf("Unexpected failure validating: %s", err)
			}

			for _, expected := range c.Expect {
				found := false
				for _, actual := range dst {
					if strings.HasPrefix(actual, expected) {
						found = true
						break
					}
				}

				if !found {
					t.Fatalf("Expected line %v not found in %v", expected, dst)
				}
			}
		})
	}
}

func TestOptionsToEnv(t *testing.T) {
	cases := []struct {
		Command *command.Command
		Args    []string
		Expect  []string
		Env     []string
	}{
		{
			Args:   []string{"test", "string", "--verbose", "--first", "something"},
			Expect: []string{"export MILPA_OPT_FIRST=something", "export MILPA_VERBOSE=true"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "string"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Options: command.Options{
					"first": &command.Option{
						Type: command.ValueTypeString,
					},
				},
			},
		},
		{
			Args:   []string{"test", "string", "--verbose", "--first", "something"},
			Expect: []string{"MILPA_OPT_FIRST=something", "MILPA_VERBOSE=true"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name: []string{"test", "string"},
					Kind: kind.Executable,
				},
				Options: command.Options{
					"first": &command.Option{
						Type: command.ValueTypeString,
					},
				},
			},
		},
		{
			Args:   []string{"test", "int", "--first", "1"},
			Expect: []string{"export MILPA_OPT_FIRST=1"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "int"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Options: command.Options{
					"first": &command.Option{
						Type: command.ValueTypeInt,
					},
				},
			},
		},
		{
			Args:   []string{"test", "bool", "--first"},
			Expect: []string{"export MILPA_OPT_FIRST=true"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "bool"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Options: command.Options{
					"first": &command.Option{
						Type: command.ValueTypeBoolean,
					},
				},
			},
		},
		{
			Args:   []string{"test", "bool-false", "--first", "false"},
			Expect: []string{"export MILPA_OPT_FIRST="},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "bool-false"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Options: command.Options{
					"first": &command.Option{
						Type: command.ValueTypeBoolean,
					},
				},
			},
		},
		{
			Args:   []string{"test", "repeated", "--pato", "quem", "--pato", "quem quem"},
			Expect: []string{"declare -a MILPA_OPT_PATO=( quem 'quem quem' )"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "repeated", "bash"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellBash,
				},
				Options: command.Options{
					"pato": &command.Option{
						Type:     command.ValueTypeString,
						Repeated: true,
					},
				},
			},
		},
		{
			Args:   []string{"test", "repeated", "--pato", "quem", "--pato", "quem quem"},
			Expect: []string{"declare -a MILPA_OPT_PATO=( quem 'quem quem' )"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "repeated", "zsh"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellZSH,
				},
				Options: command.Options{
					"pato": &command.Option{
						Type:     command.ValueTypeString,
						Repeated: true,
					},
				},
			},
		},
		{
			Args:   []string{"test", "repeated", "--pato", "quem", "--pato", "quem quem"},
			Expect: []string{"set -x MILPA_OPT_PATO quem 'quem quem'"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name:  []string{"test", "repeated", "fish"},
					Kind:  kind.ShellScript,
					Shell: kind.ShellFish,
				},
				Options: command.Options{
					"pato": &command.Option{
						Type:     command.ValueTypeString,
						Repeated: true,
					},
				},
			},
		},
		{
			Args:   []string{"test", "repeated", "--pato", "quem", "--pato", "quem quem"},
			Expect: []string{"MILPA_OPT_PATO=quem|quem quem"},
			Command: &command.Command{
				Meta: meta.Meta{
					Name: []string{"test", "repeated", "executable"},
					Kind: kind.Executable,
				},
				Options: command.Options{
					"pato": &command.Option{
						Type:     command.ValueTypeString,
						Repeated: true,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Command.FullName(), func(t *testing.T) {
			dst := []string{}
			c.Command.Path = c.Command.Meta.(meta.Meta).Name
			c.Command.SetBindings()
			fs := c.Command.FlagSet()
			fs.Bool("verbose", false, "")

			if err := fs.Parse(c.Args); err != nil {
				t.Fatalf("Could not parse test arguments (%+v): %s", c.Args, err)
			}
			c.Command.Options.Parse(c.Command.FlagSet())
			runtime.OptionsToEnv(c.Command, &dst)

			if err := c.Command.Options.AreValid(); err != nil {
				t.Fatalf("Unexpected failure validating: %s", err)
			}

			for _, expected := range c.Expect {
				found := false
				for _, actual := range dst {
					if strings.HasPrefix(actual, expected) {
						found = true
						break
					}
				}

				if !found {
					t.Fatalf("Expected line %v not found in %v", expected, dst)
				}
			}
		})
	}
}
