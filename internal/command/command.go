// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/command/runtime"
	"github.com/unrob/milpa/internal/errors"
	"gopkg.in/yaml.v3"
)

const badSpecTpl = `---
# ⚠️ Could not validate spec ⚠️

Looks like the spec for this command has errors that prevented parsing:

- command: **milpa %s**
- path:  **%s**
- error: **%s**

Run ﹅milpa itself doctor﹅ to diagnose your installed commands.

---`

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
	case kind.ShellScript, kind.Source:
		cmd.Action = runtime.Run
		spec = strings.TrimSuffix(path, filepath.Ext(path)) + ".yaml"
		logger.Main.Debugf("inner spec %s", spec)
		if m.Shell == kind.ShellUnknown {
			return cmd, fmt.Errorf("could not find a shell to run %s", path)
		}
	default:
		logger.Main.Fatalf("unknown kind: %s", m.Kind)
	}

	logger.Main.Debugf("loading spec: %s (%s)", strings.TrimSuffix(path, filepath.Ext(path))+".yaml", spec)
	var contents []byte
	if contents, err = os.ReadFile(spec); err == nil {
		err = yaml.Unmarshal(contents, cmd)
	}

	if err != nil {
		// todo: output better errors, decode yaml.TypeError
		retErr := errors.ConfigError{
			Err:    err,
			Config: spec,
			Help:   fmt.Sprintf(badSpecTpl, cmd.FullName(), spec, err),
		}
		logger.Error("failing from new")
		m.Issues = append(m.Issues, retErr)
		cmd.Meta = m
		cmd.HelpFunc = func(printLinks bool) string {
			return fmt.Sprintf(retErr.Help, err)
		}
		cmd.Action = runtime.CanRun

		return cmd, retErr
	}

	cmd.Meta = m
	return cmd.SetBindings(), nil
}
