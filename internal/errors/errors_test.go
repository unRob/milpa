// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package errors_test

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"git.rob.mx/nidito/chinampa/pkg/env"
	"github.com/unrob/milpa/internal/errors"
)

func TestSpecError(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	errs := []struct {
		Err      errors.SpecError
		Expected string
		Template string
	}{
		{
			Err: errors.SpecError{
				Path:        "config.yaml",
				CommandName: "bad template",
				Err:         fmt.Errorf("problem"),
			},
			Expected: `could not validate spec at config.yaml: problem`,
			Template: `{{ .Bad }}`,
		},
		{
			Err: errors.SpecError{
				Path:        "config.yaml",
				CommandName: "good template",
				Err:         fmt.Errorf("problem"),
			},
			Expected: `config.yaml good template problem`,
			Template: `{{ .Path }} {{ .CommandName }} {{ .Err }}`,
		},
		// 		{
		// 			Err: errors.SpecError{
		// 				Path: "config.yaml",
		// 				Err:  fmt.Errorf("problem"),
		// 			},
		// 			Expected: "Invalid configuration at config.yaml: problem",
		// 			Help:     "Invalid configuration at config.yaml: problem",
		// 		},
		// 		{
		// 			Err: errors.SpecError{
		// 				Path:        "config.yaml",
		// 				Err:         fmt.Errorf("problem"),
		// 				CommandName: "test",
		// 			},
		// 			Expected: "Invalid configuration at config.yaml: problem",
		// 			Help: `> something
		// **help**`,
		// 		},
	}

	originalTemplate := template.Must(errors.BadSpecTpl.Clone())
	defer func() {
		errors.BadSpecTpl = originalTemplate
	}()
	for _, e := range errs {
		name := e.Err.CommandName
		t.Run(name, func(t *testing.T) {
			os.Setenv(env.HelpStyle, "markdown")
			errors.BadSpecTpl = template.Must(template.New("").Parse(e.Template))

			got := e.Err.Error()
			if got != e.Expected {
				t.Fatalf("expected %s, got %s", e.Expected, got)
			}
		})
	}
}
