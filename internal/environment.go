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
package internal

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	_c "github.com/unrob/milpa/internal/constants"
)

func (cmd *Command) ToEval(args []string) string {
	output := []string{
		fmt.Sprintf("export %s=%s", _c.OutputCommandName, shellescape.Quote(cmd.FullName())),
		fmt.Sprintf("export %s=%s", _c.OutputCommandKind, shellescape.Quote(cmd.Meta.Kind)),
		fmt.Sprintf("export %s=%s", _c.OutputCommandRepo, shellescape.Quote(cmd.Meta.Repo)),
		fmt.Sprintf("export %s=%s", _c.OutputCommandPath, shellescape.Quote(cmd.Meta.Path)),
	}

	cmd.Options.ToEnv(cmd, &output)
	cmd.Arguments.ToEnv(cmd, &output)

	for idx, arg := range args {
		args[idx] = shellescape.Quote(arg)
	}
	output = append(output, "set -- "+strings.Join(args, " "))

	return strings.Join(output, "\n")
}
