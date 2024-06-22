// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package bootstrap_test

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"git.rob.mx/nidito/chinampa/pkg/logger"
	. "github.com/unrob/milpa/internal/bootstrap"
	_c "github.com/unrob/milpa/internal/constants"
	merrors "github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/repo"
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
	repo.Path = []string{}
}

func TestCheckMilpaPathSet(t *testing.T) {
	repo.Path = []string{"a", "b"}

	if err := repo.CheckPathSet(); err != nil {
		t.Fatalf("Got error with set MILPA_PATH: %v", err)
	}

	repo.Path = []string{}
	if err := repo.CheckPathSet(); err == nil {
		t.Fatalf("Got no error with unset MILPA_PATH")
	}
}

func TestBootstrapErrorsOnFakeMilpaRoot(t *testing.T) {
	resetMilpaPath()
	os.Setenv(_c.EnvVarMilpaRoot, "fake_dir")
	err := Run()
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
	err := Run()
	expected := merrors.EnvironmentError{Err: fmt.Errorf("milpa's built-in repo at %s/internal/.milpa is not a directory", root)}
	if err == nil {
		t.Fatalf("incomplete directory did not raise error, MilpaPath is %s", repo.Path)
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
		err := Run()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(repo.Path) != 1 && repo.Path[0] != root {
			t.Fatalf("Unexpected milpa path: %s", repo.Path)
		}
	})

	t.Run("MILPA_ROOT set, MILPA_PATH_PARSED unset, bad MILPA_PATH", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarMilpaPath, root+"asdfasdfasdf")
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")
		repo.Path = repo.ParsePath()

		err := Run()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(repo.Path) != 1 || repo.Path[0] != root+"/.milpa" {
			t.Fatalf("Unexpected milpa path: %s", repo.Path)
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
		repo.Path = repo.ParsePath()

		err := Run()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if !reflect.DeepEqual(repo.Path, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got: %s", expected, repo.Path)
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
		repo.Path = repo.ParsePath()
		expected := []string{root, root + "/repos/fake"}
		err := Run()

		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if !reflect.DeepEqual(repo.Path, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got: %s", expected, repo.Path)
		}
	})

	t.Run("no MILPA_ROOT", func(t *testing.T) {
		resetMilpaPath()
		os.Unsetenv(_c.EnvVarMilpaRoot)
		// update default var though, because otherwise we'd need milpa installed locally
		repo.Root = root
		err := Run()
		if err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if len(repo.Path) != 1 || repo.Path[0] != root+"/.milpa" {
			t.Fatalf("Unexpected milpa path: %s", repo.Path)
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
	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 || repo.Path[0] != root+"/.milpa" {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
	}

	// no MILPA_ROOT
	resetMilpaPath()
	os.Unsetenv(_c.EnvVarMilpaRoot)
	repo.Root = root
	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 || repo.Path[0] != root+"/.milpa" {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
	}
}

func TestBootstrapWithGit(t *testing.T) {
	root := fromProjectRoot()
	resetMilpaPath()
	os.Setenv(_c.EnvVarMilpaRoot, root)
	os.Unsetenv(_c.EnvVarLookupGitDisabled)
	os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
	os.Setenv(_c.EnvVarLookupUserReposDisabled, "true")

	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 && repo.Path[0] != root {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
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

	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 && repo.Path[0] != root {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
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

	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 2 && repo.Path[0] != wd && repo.Path[0] != wd+"/repos/internal" {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
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

	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 && repo.Path[0] != wd {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
	}

	resetMilpaPath()
	os.Setenv("XDG_DATA_HOME", "something-wrong")
	if err := Run(); err != nil {
		t.Fatalf("repo bootstrap raised unexpected error: %s", err)
	}

	if len(repo.Path) != 1 && repo.Path[0] != wd {
		t.Fatalf("Unexpected milpa path: %s", repo.Path)
	}
}

func TestBootstrapWithUserRepo(t *testing.T) {
	root := fromProjectRoot()
	home := root + "/internal/bootstrap/testdata/home"
	userRepo := home + "/.local/share/milpa/repos/user-repo"
	expected := []string{root + "/.milpa", userRepo}

	t.Run("with XDG_DATA_HOME", func(t *testing.T) {
		resetMilpaPath()
		os.Setenv(_c.EnvVarMilpaRoot, root)
		os.Setenv(_c.EnvVarLookupGitDisabled, "true")
		os.Setenv(_c.EnvVarLookupGlobalReposDisabled, "true")
		os.Unsetenv(_c.EnvVarLookupUserReposDisabled)
		os.Unsetenv("HOME")
		os.Setenv("XDG_DATA_HOME", home+"/.local/share")

		buff := &bytes.Buffer{}
		logger.SetOutput(buff)
		if err := Run(); err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if buff.String() != "" {
			t.Fatalf("repo bootstrap printed unexpected output: %s", buff)
		}

		if !reflect.DeepEqual(repo.Path, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got %s", expected, repo.Path)
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
		logger.SetOutput(buff)
		if err := Run(); err != nil {
			t.Fatalf("repo bootstrap raised unexpected error: %s", err)
		}

		if buff.String() != "" {
			t.Fatalf("repo bootstrap printed unexpected output: %s", buff)
		}

		if !reflect.DeepEqual(repo.Path, expected) {
			t.Fatalf("Unexpected milpa path: wanted %s, got %s", expected, repo.Path)
		}
	})
}
