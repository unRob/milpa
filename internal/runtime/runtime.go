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
package runtime

import (
	"fmt"
	"os"
	"strings"

	_c "github.com/unrob/milpa/internal/constants"
)

func DoctorModeEnabled() bool {
	return len(os.Args) >= 2 && (os.Args[1] == "__doctor" || (len(os.Args) > 2 && (os.Args[1] == "itself" && os.Args[2] == "doctor")))
}

func ValidationEnabled() bool {
	return os.Getenv(_c.EnvVarValidationDisabled) != "1"
}

func VerboseEnabled() bool {
	return os.Getenv(_c.EnvVarMilpaVerbose) != ""
}

func ColorEnabled() bool {
	return os.Getenv(_c.EnvVarMilpaUnstyled) == "" && os.Getenv(_c.EnvVarHelpUnstyled) == ""
}

func UnstyledHelpEnabled() bool {
	return os.Getenv(_c.EnvVarHelpUnstyled) == "enabled"
}

var MilpaPath = strings.Split(os.Getenv(_c.EnvVarMilpaPath), ":")

func CheckMilpaPathSet() error {
	if len(MilpaPath) == 0 {
		return fmt.Errorf("no %s set on the environment", _c.EnvVarMilpaPath)
	}
	return nil
}
