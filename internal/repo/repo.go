// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package repo

import (
	"fmt"
	"os"
	"strings"

	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
)

// Path holds a list of paths where milpa looks for repos.
var Path = ParsePath()

// ParsePath turns MILPA_PATH into a string slice.
func ParsePath() []string {
	return strings.Split(os.Getenv(_c.EnvVarMilpaPath), ":")
}

func CheckPathSet() error {
	if len(Path) == 0 {
		return errors.EnvironmentError{
			Err: fmt.Errorf("no %s set on the environment", _c.EnvVarMilpaPath),
		}
	}
	return nil
}

// Root points to the system's milpa installation.
var Root = "/usr/local/lib/milpa"
