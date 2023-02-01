// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/unrob/milpa/internal/bootstrap"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"gopkg.in/yaml.v3"
)

func posixSource(executable string, cmd *command.Command, meta Meta) error {
	env := ToEval(cmd)

	out, err := os.CreateTemp(os.TempDir(), "milpa-cmdenv-*")
	if err != nil {
		return err
	}

	_, err = out.Write([]byte(env))
	if err != nil {
		return fmt.Errorf("could not write to temporary file: %s", err)
	}

	itself, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not tell our executable path: %s", err)
	}
	cmdEnv := append(os.Environ(),
		_c.EnvVarMilpaRoot+"="+bootstrap.MilpaRoot,
		_c.OutputCommandPath+"="+meta.Path,
		"MILPA="+itself,
	)

	beforeHook := meta.Repo + "/hooks/before-run.sh"
	sources := strings.Join([]string{
		"source '" + out.Name() + "'",
		"source '" + bootstrap.MilpaRoot + "/.milpa/utils.sh'",
		"[[ -f '" + beforeHook + "' ]] && source '" + beforeHook + "'",
	}, ";") + ";"

	args := []string{
		executable,
		"-c",
		"set -o allexport;" + sources + "set +o allexport; rm " + out.Name() + "; source " + meta.Path + ";",
	}

	logrus.Debugf("calling %s", args)

	return syscall.Exec(executable, args, cmdEnv)
}

func New(path string, repo string) (cmd *command.Command, err error) {
	meta := metaForPath(path, repo)
	cmd = &command.Command{
		Path:      meta.Name,
		Arguments: []*command.Argument{},
		Options:   command.Options{},
	}

	var spec string
	switch meta.Kind {
	case KindVirtual:
		spec = path
	case KindExecutable:
		spec = path + ".yaml"
	case KindPosix, KindSource:
		spec = strings.TrimSuffix(path, filepath.Ext(path)) + ".yaml"
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
		meta.issues = append(meta.issues, err)
		cmd.Meta = meta
		cmd.HelpFunc = func(printLinks bool) string {
			return `---
# ⚠️ Could not validate spec ⚠️

Looks like the spec for this command has errors that prevented parsing:

**` + fmt.Sprint(err) + `**

Run ﹅milpa itself doctor﹅ to diagnose your installed commands.

---`
		}
		cmd.Action = canRun

		return cmd, err
	}

	switch meta.Kind {
	case KindPosix:
		cmd.Action = func(cmd *command.Command) error {
			if err := canRun(cmd); err != nil {
				return err
			}
			logger.Main.Debugf("running command")

			if meta.Shell == "" {
				return fmt.Errorf("could not find a shell to run %s", path)
			}

			shell, err := exec.LookPath(meta.Shell)
			if err != nil {
				return fmt.Errorf("could not find an executable for %s: %s", shell, err)
			}
			return posixSource(shell, cmd, meta)
		}
	case KindExecutable:
		cmd.Action = func(cmd *command.Command) error {
			if err := canRun(cmd); err != nil {
				return err
			}
			logger.Main.Debugf("running command")

			cmdEnv := Env(cmd, os.Environ())
			args := ArgumentsToSlice(cmd)

			// Launch command with user provided arguments
			return syscall.Exec(meta.Path, args, cmdEnv) // nolint:gosec
		}
	}

	cmd.Meta = meta
	return cmd.SetBindings(), nil
}

func canRun(cmd *command.Command) error {
	if cmd.Meta == nil {
		return fmt.Errorf("unknown meta: %s", cmd.Path)
	}

	meta := cmd.Meta.(Meta)

	if len(meta.issues) > 0 {
		issues := []string{}
		for _, i := range meta.issues {
			issues = append(issues, i.Error())
		}

		return errors.ConfigError{
			Err: fmt.Errorf("cannot run command <%s>: %s", cmd.FullName(), strings.Join(issues, "\n")),
		}
	}
	return nil
}
