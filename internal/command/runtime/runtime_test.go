// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime_test

import (
	"fmt"
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/command/runtime"
)

func TestCanRun(t *testing.T) {
	cases := []struct {
		Command  *command.Command
		Expected string
	}{
		{
			Command: &command.Command{
				Path: []string{"no-meta"},
			},
			Expected: "unknown meta for command: no-meta",
		},
		{
			Command: &command.Command{
				Path: []string{"bad-meta"},
				Meta: 42,
			},
			Expected: "unknown meta type for command bad-meta: int",
		},
		{
			Command: &command.Command{
				Path: []string{"bad-meta"},
				Meta: meta.Meta{
					Error: fmt.Errorf("problem"),
				},
			},
			Expected: "problem",
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.Command.Path, "-"), func(t *testing.T) {
			err := runtime.CanRun(tc.Command)
			if tc.Expected != "" && err == nil || err.Error() != tc.Expected {
				t.Errorf("expected error %s, got %s", tc.Expected, err)
			}
		})
	}
}
