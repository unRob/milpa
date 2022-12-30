// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package errors

import "fmt"

type NotFound struct {
	Msg   string
	Group []string
}

type BadArguments struct {
	Msg string
}

type NotExecutable struct {
	Msg string
}

type ConfigError struct {
	Err    error
	Config string
}

type EnvironmentError struct {
	Err error
}

type SubCommandExit struct {
	Err      error
	ExitCode int
}

func (err NotFound) Error() string {
	return err.Msg
}

func (err BadArguments) Error() string {
	return err.Msg
}

func (err NotExecutable) Error() string {
	return err.Msg
}

func (err SubCommandExit) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}

	return ""
}

func (err ConfigError) Error() string {
	return fmt.Sprintf("Invalid configuration %s: %v", err.Config, err.Err)
}

func (err EnvironmentError) Error() string {
	return fmt.Sprintf("Invalid MILPA_ environment: %v", err.Err)
}
