// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions_test

import (
	"bytes"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/env"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "github.com/unrob/milpa/internal/actions"
	"github.com/unrob/milpa/internal/bootstrap"
	"github.com/unrob/milpa/internal/lookup"
)

func fromProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
	wd, _ := os.Getwd()
	return wd
}

func TestRenderDocsHelpMain(t *testing.T) {
	lookup.DefaultFS = os.DirFS("/")
	root := fromProjectRoot()
	bootstrap.MilpaPath = []string{root + "/.milpa"}
	AfterHelp = func(code int) {}
	Docs.SetBindings()
	out := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd := &cobra.Command{Use: "asdf"}
	os.Setenv(env.HelpStyle, "markdown")
	cmd.SetHelpFunc(Docs.HelpRenderer(command.Root.Options))
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	Docs.Cobra = cmd
	logrus.SetLevel(logrus.DebugLevel)
	err := Docs.Run(cmd, []string{})
	if err != nil {
		t.Fatalf("did not run: %s", out.String())
	}

	expected := `
## Available topics:

- [milpa](milpa)
`
	if !strings.Contains(out.String(), expected) {
		t.Fatalf("Did not find <%s> in output, got: \n\n%s", expected, out.String())
	}
}

func TestRenderDocsRender(t *testing.T) {
	lookup.DefaultFS = os.DirFS("/")
	root := fromProjectRoot()
	bootstrap.MilpaPath = []string{root + "/.milpa"}
	AfterHelp = func(code int) {}
	Docs.SetBindings()
	out := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd := &cobra.Command{Use: "asdf"}
	os.Setenv(env.HelpStyle, "markdown")
	cmd.SetHelpFunc(Docs.HelpRenderer(command.Root.Options))
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	Docs.Cobra = cmd
	logrus.SetLevel(logrus.DebugLevel)
	err := Docs.Run(cmd, []string{"milpa"})
	if err != nil {
		t.Fatalf("did not run: %s", err)
	}

	expected := `# milpa

[` + "`milpa`" + `](https://milpa.dev) is a command-line tool to care for one's own garden of scripts. [Its name](https://en.wikipedia.org/wiki/Milpa) comes from an agricultural method that combines multiple crops in close proximity.`
	if !strings.HasPrefix(out.String(), expected) {
		t.Fatalf("Did not find <%s> in output, got: \n\n%s", expected, out.String())
	}
}
