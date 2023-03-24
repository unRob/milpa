// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package command

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"gopkg.in/yaml.v3"
)

func New(path string, repo string) (cmd *command.Command, err error) {
	meta := metaForPath(path, repo)
	cmd = &command.Command{
		Path: meta.Name,
		Action: func(cmd *command.Command) error {
			if err := canRun(cmd); err != nil {
				return err
			}
			logger.Main.Debugf("running command")

			env := ToEval(cmd, []string{})

			if os.Getenv(_c.EnvVarCompaOut) != "" {
				return os.WriteFile(os.Getenv(_c.EnvVarCompaOut), []byte(env), 0600)
			}

			fmt.Println(env)
			return nil
		},
	}
	cmd.Arguments = []*command.Argument{}
	cmd.Options = command.Options{}

	spec := strings.TrimSuffix(path, ".sh") + ".yaml"
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
