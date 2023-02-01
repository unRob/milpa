// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package util

import (
	"os"
	"strconv"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/env"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	_c "github.com/unrob/milpa/internal/constants"
)

var trueIshValues = []string{
	"1",
	"yes",
	"true",
	"enable",
	"enabled",
	"on",
	"always",
}

func IsTrueIsh(val string) bool {
	for _, positive := range trueIshValues {
		if val == positive {
			return true
		}
	}

	return false
}

// EnvironmentMap returns the resolved environment map.
func EnvironmentMap(mp []string, mr string) map[string]string {
	res := map[string]string{
		_c.EnvVarMilpaRoot:       mr,
		_c.EnvVarMilpaPath:       strings.Join(mp, ":"),
		_c.EnvVarMilpaPathParsed: "true",
	}
	trueString := strconv.FormatBool(true)

	if !runtime.ColorEnabled() {
		res[env.NoColor] = trueString
	} else if IsTrueIsh(os.Getenv(env.ForceColor)) {
		res[env.ForceColor] = "always"
	}

	if runtime.DebugEnabled() {
		res[env.Debug] = trueString
	}

	if runtime.VerboseEnabled() {
		res[env.Verbose] = trueString
	} else if IsTrueIsh(os.Getenv(env.Silent)) {
		res[env.Silent] = trueString
	}

	return res
}
