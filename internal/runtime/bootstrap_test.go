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
package runtime_test

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	_c "github.com/unrob/milpa/internal/constants"
	merrors "github.com/unrob/milpa/internal/errors"
	. "github.com/unrob/milpa/internal/runtime"
)

func fromProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
	wd, _ := os.Getwd()
	return wd
}

func resetMilpaPath() {
	os.Unsetenv(_c.EnvVarMilpaPathParsed)
	os.Unsetenv(_c.EnvVarMilpaPath)
	MilpaPath = []string{}
}

func TestBootstrapErrorsOnFakeMilpaRoot(t *testing.T) {
	resetMilpaPath()
	os.Setenv(_c.EnvVarMilpaRoot, "fake_dir")
	err := Bootstrap()
	expected := merrors.EnvironmentError{Err: fmt.Errorf("MILPA_ROOT (fake_dir) is not a directory")}
	if err == nil {
		t.Fatal("fake directory did not raise error")
	}

	if cErr, ok := err.(merrors.EnvironmentError); !ok {
		t.Fatalf("bootstrap did not fail with correct error.\nexpected %s, got %s", expected, err)
	} else if cErr.Error() != expected.Error() {
		t.Fatalf("bootstrap did not fail with correct error message.\nwant %s\nhave %s", expected.Error(), err.Error())
	}
}

func TestBootstrapErrorsOnIncompleteMilpaRoot(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()

	// this is a real directory, but without a .milpa dir!
	os.Setenv(_c.EnvVarMilpaRoot, root+"/internal")
	os.Setenv(_c.EnvVarLookupGitDisabled, "true")
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")
	err := Bootstrap()
	expected := merrors.EnvironmentError{Err: fmt.Errorf("milpa's built-in repo at %s/internal/.milpa is not a directory", root)}
	if err == nil {
		t.Fatalf("incomplete directory did not raise error, MilpaPath is %s", MilpaPath)
	}

	if cErr, ok := err.(merrors.EnvironmentError); !ok {
		t.Fatalf("incomplete bootstrap did not fail with correct error.\nexpected %s, got %s", expected, err)
	} else if cErr.Error() != expected.Error() {
		t.Fatalf("incomplete bootstrap did not fail with correct error message.\nwant %s\nhave %s", expected.Error(), err.Error())
	}
}

func TestBootstrapWithMilpaPath(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()

	t.Run("empty string", func(t *testing.T) {
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarMilpaPath, "")
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

		// with MILPA_ROOT set
		err := Bootstrap()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(MilpaPath) != 1 && MilpaPath[0] != root {
			t.Fatalf("Unexpected milpa path: %s", MilpaPath)
		}
	})

	t.Run("MILPA_ROOT set, MILPA_PATH_PARSED unset, bad MILPA_PATH", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarMilpaPath, root+"asdfasdfasdf")
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")
		MilpaPath = ParseMilpaPath()

		err := Bootstrap()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(MilpaPath) != 1 || MilpaPath[0] != root+"/.milpa" {
			t.Fatalf("Unexpected milpa path: %s", MilpaPath)
		}
	})

	t.Run("MILPA_ROOT set, MILPA_PATH_PARSED unset, correct MILPA_PATH", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarMilpaPath, root+"/repos/internal")
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

		expected := []string{root + "/repos/internal/.milpa", root + "/.milpa"}
		MilpaPath = ParseMilpaPath()

		err := Bootstrap()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if !reflect.DeepEqual(MilpaPath, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got: %s", expected, MilpaPath)
		}
	})

	t.Run("MILPA_ROOT set, MILPA_PATH_PARSED set", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarMilpaPath, strings.Join([]string{root, root + "/repos/fake"}, ":"))
		os.Setenv(_c.EnvVarMilpaPathParsed, "true")
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")
		MilpaPath = ParseMilpaPath()
		expected := []string{root, root + "/repos/fake"}
		err := Bootstrap()

		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if !reflect.DeepEqual(MilpaPath, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got: %s", expected, MilpaPath)
		}
	})

	t.Run("no MILPA_ROOT", func(t *testing.T) {
		resetMilpaPath()
		os.Unsetenv(_c.EnvVarMilpaRoot)
		// update default var though, because otherwise we'd need milpa installed locally
		MilpaRoot = root
		err := Bootstrap()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(MilpaPath) != 1 || MilpaPath[0] != root+"/.milpa" {
			t.Fatalf("Unexpected milpa path: %s", MilpaPath)
		}
	})
}

func TestBootstrapOkOnRepo(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()
	os.Setenv(_c.EnvVarMilpaRoot, root)
	os.Setenv(_c.EnvVarLookupGitDisabled, "true")
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

	// with MILPA_ROOT set
	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 || MilpaPath[0] != root+"/.milpa" {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}

	// no MILPA_ROOT
	resetMilpaPath()
	os.Unsetenv(_c.EnvVarMilpaRoot)
	MilpaRoot = root
	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 || MilpaPath[0] != root+"/.milpa" {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}
}

func TestBootstrapWithGit(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()
	os.Setenv(_c.EnvVarMilpaRoot, root)
	os.Unsetenv(_c.EnvVarLookupGitDisabled)
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 && MilpaPath[0] != root {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}
}

func TestBootstrapWithoutGit(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()
	ospath := os.Getenv("PATH")
	os.Setenv("PATH", "")

	defer func() {
		os.Setenv("PATH", ospath)
	}()

	os.Setenv(_c.EnvVarMilpaRoot, root)
	os.Unsetenv(_c.EnvVarLookupGitDisabled)
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 && MilpaPath[0] != root {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}
}

func TestBootstrapWithGlobalRepo(t *testing.T) {
	fromProjectRoot()
	resetMilpaPath()
	wd, _ := os.Getwd()
	os.Setenv(_c.EnvVarMilpaRoot, wd)
	os.Setenv(_c.EnvVarLookupGitDisabled, "true")
	os.Unsetenv(_c.EnvVarLookupGlobalReposDisabled)
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 2 && MilpaPath[0] != wd && MilpaPath[0] != wd+"/repos/internal" {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}
}

func TestBootstrapWithUserRepoAndNoHome(t *testing.T) {
	fromProjectRoot()
	resetMilpaPath()
	wd, _ := os.Getwd()
	os.Setenv(_c.EnvVarMilpaRoot, wd)
	os.Setenv(_c.EnvVarLookupGitDisabled, "true")
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Unsetenv(_c.EnvVarLookupUserReposDisabled)

	os.Unsetenv("XDG_DATA_HOME")
	os.Unsetenv("HOME")

	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 && MilpaPath[0] != wd {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}

	resetMilpaPath()
	os.Setenv("XDG_DATA_HOME", "something-wrong")
	if err := Bootstrap(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(MilpaPath) != 1 && MilpaPath[0] != wd {
		t.Fatalf("Unexpected milpa path: %s", MilpaPath)
	}
}

func TestBootstrapWithUserRepo(t *testing.T) {
	root := fromProjectRoot()
	home := root + "/internal/runtime/testdata/home"
	repo := home + "/.local/share/milpa/repos/user-repo"
	expected := []string{root + "/.milpa", repo}

	t.Run("with XDG_DATA_HOME", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Unsetenv(_c.EnvVarLookupUserReposDisabled)
		os.Unsetenv("HOME")
		os.Setenv("XDG_DATA_HOME", home)

		buff := &bytes.Buffer{}
		logrus.SetOutput(buff)
		if err := Bootstrap(); err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if buff.String() != "" {
			t.Fatalf("repo bootstrap printed unexpected output: %s", buff)
		}

		if !reflect.DeepEqual(MilpaPath, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got %s", expected, MilpaPath)
		}
	})

	t.Run("with HOME", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Unsetenv(_c.EnvVarLookupUserReposDisabled)
		os.Unsetenv("XDG_DATA_HOME")
		os.Setenv("HOME", home)

		buff := &bytes.Buffer{}
		logrus.SetOutput(buff)
		if err := Bootstrap(); err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if buff.String() != "" {
			t.Fatalf("repo bootstrap printed unexpected output: %s", buff)
		}

		if !reflect.DeepEqual(MilpaPath, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got %s", expected, MilpaPath)
		}
	})
}
