// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/command/runtime"
	"github.com/unrob/milpa/internal/errors"
	"gopkg.in/yaml.v3"
)

func New(path string, repo string) (cmd *command.Command, err error) {
	m := meta.ForPath(path, repo)
	cmd = &command.Command{
		Path:      m.Name,
		Arguments: []*command.Argument{},
		Options:   command.Options{},
	}

	var spec string
	switch m.Kind {
	case kind.Virtual:
		spec = path
	case kind.Executable:
		cmd.Action = runtime.Run
		spec = path + ".yaml"
	case kind.Posix, kind.Source:
		cmd.Action = runtime.Run
		spec = strings.TrimSuffix(path, filepath.Ext(path)) + ".yaml"
		if m.Shell == "" {
			return cmd, fmt.Errorf("could not find a shell to run %s", path)
		}
	}

	var contents []byte
	if contents, err = os.ReadFile(spec); err == nil {
		err = yaml.Unmarshal(contents, cmd)
	}

	if err != nil {
		// todo: output better errors, decode yaml.TypeError
		err = errors.ConfigError{
			Err:    err,
			Config: spec,
		}
		m.Issues = append(m.Issues, err)
		cmd.Meta = m
		cmd.HelpFunc = func(printLinks bool) string {
			return `---
# ⚠️ Could not validate spec ⚠️

Looks like the spec for this command has errors that prevented parsing:

**` + fmt.Sprint(err) + `**

Run ﹅milpa itself doctor﹅ to diagnose your installed commands.

---`
		}
		cmd.Action = runtime.CanRun

		return cmd, err
	}

	cmd.Meta = m
	return cmd.SetBindings(), nil
}
