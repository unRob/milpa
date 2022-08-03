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
package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/unrob/milpa/internal/actions"
	"github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/registry"
	"github.com/unrob/milpa/internal/runtime"
)

var version = "beta"

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceColors:            runtime.ColorEnabled(),
	})

	root := actions.RootCommand(version)
	f := root.PersistentFlags()
	silent := false
	if err := f.Parse(os.Args); err != nil {
		silent, _ = f.GetBool("silent")
	}

	if !silent && runtime.DebugEnabled() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	isDoctor := runtime.DoctorModeEnabled()
	logrus.Debugf("doctor mode enabled: %v", isDoctor)

	err := runtime.Bootstrap()
	if err != nil {
		logrus.Fatal(err)
	}

	err = registry.FindAllSubCommands(!isDoctor)
	if err != nil && !isDoctor {
		logrus.Fatalf("Could not find subcommands: %s", err)
	}

	registry.SetRoot(root, actions.Root)
	initialArgs := []string{"milpa"}
	os.Args = append(initialArgs, os.Args[1:]...) //nolint:gocritic

	errors.HandleCobraExit(root.ExecuteC())
}
