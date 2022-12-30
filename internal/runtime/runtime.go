// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_c "github.com/unrob/milpa/internal/constants"
)

var MilpaPath = ParseMilpaPath()

// ParseMilpaPath turns MILPA_PATH into a string slice.
func ParseMilpaPath() []string {
	return strings.Split(os.Getenv(_c.EnvVarMilpaPath), ":")
}

var falseIshValues = []string{
	"",
	"0",
	"no",
	"false",
	"disable",
	"disabled",
	"off",
	"never",
}

var trueIshValues = []string{
	"1",
	"yes",
	"true",
	"enable",
	"enabled",
	"on",
	"always",
}

func isFalseIsh(val string) bool {
	for _, negative := range falseIshValues {
		if val == negative {
			return true
		}
	}

	return false
}

func isTrueIsh(val string) bool {
	for _, positive := range trueIshValues {
		if val == positive {
			return true
		}
	}

	return false
}

func DoctorModeEnabled() bool {
	count := len(os.Args)
	if count < 2 {
		return false
	}
	first := os.Args[1]

	return first == "__doctor" || count >= 2 && (first == "itself" && os.Args[2] == "doctor")
}

func DebugEnabled() bool {
	return isTrueIsh(os.Getenv(_c.EnvVarDebug))
}

func ValidationEnabled() bool {
	return isFalseIsh(os.Getenv(_c.EnvVarValidationDisabled))
}

func VerboseEnabled() bool {
	return isTrueIsh(os.Getenv(_c.EnvVarMilpaVerbose))
}

func ColorEnabled() bool {
	return isFalseIsh(os.Getenv(_c.EnvVarMilpaUnstyled)) && !UnstyledHelpEnabled()
}

func UnstyledHelpEnabled() bool {
	return isTrueIsh(os.Getenv(_c.EnvVarHelpUnstyled))
}

func CheckMilpaPathSet() error {
	if len(MilpaPath) == 0 {
		return fmt.Errorf("no %s set on the environment", _c.EnvVarMilpaPath)
	}
	return nil
}

// EnvironmentMap returns the resolved environment map.
func EnvironmentMap() map[string]string {
	env := map[string]string{}
	env[_c.EnvVarMilpaPath] = strings.Join(MilpaPath, ":")
	trueString := strconv.FormatBool(true)
	env[_c.EnvVarMilpaPathParsed] = trueString

	if !ColorEnabled() {
		env[_c.EnvVarMilpaUnstyled] = trueString
	} else if isTrueIsh(os.Getenv(_c.EnvVarMilpaForceColor)) {
		env[_c.EnvVarMilpaForceColor] = "always"
	}

	if DebugEnabled() {
		env[_c.EnvVarDebug] = trueString
	}

	if VerboseEnabled() {
		env[_c.EnvVarMilpaVerbose] = trueString
	} else if isTrueIsh(os.Getenv(_c.EnvVarMilpaSilent)) {
		env[_c.EnvVarMilpaSilent] = trueString
	}

	return env
}
