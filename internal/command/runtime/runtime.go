// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/unrob/milpa/internal/bootstrap"
	"github.com/unrob/milpa/internal/command/kind"
	"github.com/unrob/milpa/internal/command/meta"
	"github.com/unrob/milpa/internal/errors"
)

var log = logger.Sub("runtime")

var posixEntrypoint = template.Must(template.New("").Parse(`set -o allexport;
source "{{ .env }}";
set +o allexport;
rm "{{ .env }}";
source "{{ .milpaRoot }}/.milpa/runtime.{{ .shell }}";
[[ -f "{{ .repo }}//hooks/before-run.sh" ]] && source "{{ .repo }}/hooks/before-run.sh";
source "{{ .path }}";`))

var fishEntrypoint = template.Must(template.New("").Parse(`source "{{ .env }}"
rm "{{ .env }}"
source "{{ .milpaRoot }}/.milpa/runtime.fish";
if test -f "{{ .repo }}/hooks/before-run.fish"
  source "{{ .repo }}/hooks/before-run.fish"
end
source "{{ .path }}"`))

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
		cErr, ok := meta.Issues[0].(errors.ConfigError)
		if len(meta.Issues) == 1 && ok {
			return cErr
		}
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
	case kind.ShellScript, kind.Source:
		return Shell(cmd)
	}

	return fmt.Errorf("no runtime available for milpa command %s with kind %s", cmd.FullName(), m.Kind)
}

// Shell replaces the current process with a shell invocation for a command.
func Shell(cmd *command.Command) error {
	m := cmd.Meta.(meta.Meta)
	shell, err := m.Shell.Path()
	if err != nil {
		return err
	}

	cmdEnv := BaseEnv(m)

	args := []string{}
	switch m.Shell {
	case kind.ShellBash, kind.ShellZSH:
		env := ToEval(cmd)

		out, err := os.CreateTemp(os.TempDir(), "milpa-cmdenv-*")
		if err != nil {
			return err
		}

		_, err = out.Write([]byte(env))
		if err != nil {
			return fmt.Errorf("could not write to temporary file: %s", err)
		}

		buf := &bytes.Buffer{}
		err = template.Must(posixEntrypoint.Clone()).Execute(buf, map[string]string{
			"env":       out.Name(),
			"repo":      m.Repo,
			"milpaRoot": bootstrap.MilpaRoot,
			"path":      m.Path,
			"shell":     string(m.Shell),
		})
		if err != nil {
			return err
		}

		rt := buf.String()
		log.Tracef("bash runtime: %s", rt)

		args = []string{shell, "-c", rt}
	case kind.ShellFish:
		env := ToEval(cmd)

		out, err := os.CreateTemp(os.TempDir(), "milpa-cmdenv-*")
		if err != nil {
			return err
		}

		_, err = out.Write([]byte(env))
		if err != nil {
			return fmt.Errorf("could not write to temporary file: %s", err)
		}
		buf := &bytes.Buffer{}
		err = template.Must(fishEntrypoint.Clone()).Execute(buf, map[string]string{
			"env":       out.Name(),
			"repo":      m.Repo,
			"milpaRoot": bootstrap.MilpaRoot,
			"path":      m.Path,
			"shell":     string(m.Shell),
		})
		if err != nil {
			return err
		}

		args = []string{shell, "-c", buf.String()}
	}

	log.Debugf("calling shell %s (%s) with command %s", m.Shell, shell, args)

	return fork(shell, args, cmdEnv)
}

// Executable replaces the current process with the forked command.
func Executable(cmd *command.Command) error {
	m := cmd.Meta.(meta.Meta)

	cmdEnv, args := Env(cmd, BaseEnv(m))

	log.Debugf("calling executable command %s", args)
	// Launch command with user provided arguments
	return fork(m.Path, args, cmdEnv) // nolint:gosec
}
