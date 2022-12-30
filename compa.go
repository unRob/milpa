// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
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
