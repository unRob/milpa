// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
