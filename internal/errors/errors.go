// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package errors

import (
	"bytes"
	"fmt"
	"text/template"

	"git.rob.mx/nidito/chinampa/pkg/logger"
	"git.rob.mx/nidito/chinampa/pkg/render"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
)

var BadSpecTpl = template.Must(template.New("").Parse(`---
# ⚠️ Could not validate spec ⚠️

Looks like the spec at **{{ .Path }}** for the command **milpa {{ .CommandName }}** has errors that prevented parsing:

**{{ .Err }}**

Run ﹅milpa itself doctor﹅ to diagnose your installed commands.

---`))

// ProgrammerError happens when either I or the repo developer did something to upset milpa.
type ProgrammerError struct {
	Err error
}

func (err ProgrammerError) Error() string {
	return err.Err.Error()
}

// SpecError happens when a command's spec is somehow invalid.
type SpecError struct {
	// Err is the upstream error that prevented parsing.
	Err error
	// Path is the filesystem path to the failing spec.
	Path string
	// CommandName is the name of the command with the broken spec.
	CommandName string
}

func (err SpecError) Doctor() error {
	return fmt.Errorf("could not validate spec at %s: %s", err.Path, err.Err.Error())
}

func (err SpecError) Error() string {
	buf := bytes.Buffer{}
	if e := template.Must(BadSpecTpl.Clone()).Execute(&buf, err); e == nil {
		if text, e := render.Markdown(buf.Bytes(), runtime.ColorEnabled()); e == nil {
			return string(text)
		}
		logger.Errorf("help render failed: %s", e)
	} else {
		logger.Errorf("could not interpolate template: %s", e)
	}

	return err.Doctor().Error()
}

// EnvironmentError happens when the environment variables and folders expected by milpa are off.
type EnvironmentError struct {
	Err error
}

func (err EnvironmentError) Error() string {
	return fmt.Sprintf("Invalid MILPA_ environment: %v", err.Err)
}
