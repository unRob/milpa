// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command_test

import (
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/command"
	. "github.com/unrob/milpa/internal/command"
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
				Meta: Meta{
					Name: []string{"test", "required", "present"},
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
				Meta: Meta{
					Name: []string{"test", "default", "present"},
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
				"declare -a MILPA_ARG_VARIADIC=(one two three)",
			},
			Command: &command.Command{
				Meta: Meta{
					Name: []string{"test", "variadic"},
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
				"declare -a MILPA_ARG_VARIADIC=('one and stuff' two three)",
			},
			Command: &command.Command{
				Meta: Meta{
					Name: []string{"test", "variadic"},
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
				Meta: Meta{
					Name: []string{"test", "static", "default"},
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
				Meta: Meta{
					Name: []string{"test", "static", "good"},
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
				Meta: Meta{
					Name: []string{"test", "script", "good"},
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
			ArgumentsToEnv(c.Command, &dst, "export ")

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
			Args:   []string{"test", "string", "--first", "something"},
			Expect: []string{"export MILPA_OPT_FIRST=something"},
			Command: &command.Command{
				Meta: Meta{
					Name: []string{"test", "string"},
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
				Meta: Meta{
					Name: []string{"test", "int"},
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
				Meta: Meta{
					Name: []string{"test", "bool"},
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
				Meta: Meta{
					Name: []string{"test", "bool-false"},
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
			Expect: []string{"export MILPA_OPT_PATO=( quem 'quem quem' )"},
			Command: &command.Command{
				Meta: Meta{
					Name: []string{"test", "repeated"},
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
			c.Command.SetBindings()
			err := c.Command.FlagSet().Parse(c.Args)
			if err != nil {
				t.Fatalf("Could not parse test arguments (%+v): %s", c.Args, err)
			}
			c.Command.Options.Parse(c.Command.FlagSet())
			OptionsToEnv(c.Command, &dst, "export ")

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
