// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package main

import (
	"os"

	"git.rob.mx/nidito/chinampa"
	"git.rob.mx/nidito/chinampa/pkg/runtime"
	"github.com/sirupsen/logrus"
	"github.com/unrob/milpa/internal/actions"
	"github.com/unrob/milpa/internal/bootstrap"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/lookup"
)

var version = "beta"

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceColors:            runtime.ColorEnabled(),
	})

	if !runtime.SilenceEnabled() && runtime.DebugEnabled() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	isDoctor := actions.DoctorModeEnabled()
	logrus.Debugf("doctor mode enabled: %v", isDoctor)

	err := bootstrap.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	cfg := chinampa.Config{
		Name:    "milpa",
		Version: version,
		Summary: "Runs commands found in " + _c.RepoRoot + " folders",
		Description: `﹅milpa﹅ is a command-line tool to care for one's own garden of scripts, its name comes from an agricultural method that combines multiple crops in close proximity. You and your team write scripts and a little spec for each command -use bash, or any other language-, and ﹅milpa﹅ provides autocompletions, sub-commands, argument parsing and validation so you can skip the toil and focus on your scripts.

See [﹅milpa help docs milpa﹅](/.milpa/docs/milpa/index.md) for more information about ﹅milpa﹅`,
	}

	chinampa.Register(actions.Doctor)
	chinampa.Register(actions.Docs)
	chinampa.Register(actions.GenerateDocs)
	chinampa.Register(actions.CommandTree)
	chinampa.Register(actions.Fetch)

	err = lookup.AllSubCommands(!isDoctor)
	if err != nil && !isDoctor {
		logrus.Fatalf("Could not find subcommands: %s", err)
	}

	if err := chinampa.Execute(cfg); err != nil {
		logrus.Errorf("total failure: %s", err)
		os.Exit(2)
	}
}
