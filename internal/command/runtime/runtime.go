// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/unrob/milpa/internal/bootstrap"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/errors"
)

var log = logger.Sub("runtime")

// CanRun is the last runtime check before actually calling a command.
func CanRun(cmd *command.Command) error {
	if cmd.Meta == nil {
		return errors.ProgrammerError{
			Err: fmt.Errorf("unknown meta: %s", cmd.Path),
		}
	}

	meta, ok := cmd.Meta.(meta.Meta)
	if !ok {
		return errors.ProgrammerError{
			Err: fmt.Errorf("meta of unknown kind for %s: %+v", cmd.Path, cmd.Meta),
		}
	}

	if len(meta.Issues) > 0 {
		issues := []string{}
		for _, i := range meta.Issues {
			issues = append(issues, i.Error())
		}

		return errors.ConfigError{
			Err: fmt.Errorf("cannot run command <%s>: %s", cmd.FullName(), strings.Join(issues, "\n")),
		}
	}
	return nil
}

// Run is the chinampa action to be take on a valid command.
func Run(cmd *command.Command) error {
	m := cmd.Meta.(meta.Meta)
	if err := CanRun(cmd); err != nil {
		return err
	}
	// logger.Main.Debugf("running command")
	switch m.Kind {
	case kind.Executable:
		return Executable(cmd)
	case kind.Posix, kind.Source:
		return Shell(cmd)
	}

	return fmt.Errorf("no runtime available for milpa command %s with kind %s", cmd.FullName(), m.Kind)
}

// Shell replaces the current process with a shell invocation for a command.
func Shell(cmd *command.Command) error {
	m := cmd.Meta.(meta.Meta)
	shell, err := exec.LookPath(m.Shell)
	if err != nil {
		return fmt.Errorf("could not find an executable for %s: %s", m.Shell, err)
	}

	env := ToEval(cmd)

	out, err := os.CreateTemp(os.TempDir(), "milpa-cmdenv-*")
	if err != nil {
		return err
	}

	_, err = out.Write([]byte(env))
	if err != nil {
		return fmt.Errorf("could not write to temporary file: %s", err)
	}

	cmdEnv := BaseEnv(m)

	beforeHook := m.Repo + "/hooks/before-run.sh"
	sources := strings.Join([]string{
		"source '" + out.Name() + "'",
		"source '" + bootstrap.MilpaRoot + "/.milpa/utils.sh'",
		"[[ -f '" + beforeHook + "' ]] && source '" + beforeHook + "'",
	}, ";") + ";"

	args := []string{
		shell,
		"-c",
		"set -o allexport;" + sources + "set +o allexport; rm " + out.Name() + "; source " + m.Path + ";",
	}

	log.Debugf("calling shell command %s", args)

	return fork(shell, args, cmdEnv)
}

// Executable replaces the current process with the forked command.
func Executable(cmd *command.Command) error {
	m := cmd.Meta.(meta.Meta)

	cmdEnv := Env(cmd, BaseEnv(m))
	args := ArgumentsToSlice(cmd)

	log.Debugf("calling executable command %s", args)
	// Launch command with user provided arguments
	return fork(m.Path, args, cmdEnv) // nolint:gosec
}
