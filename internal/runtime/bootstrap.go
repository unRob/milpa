// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package runtime

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
)

// MilpaRoot points to the system's milpa installation.
var MilpaRoot = "/usr/local/lib/milpa"

func Bootstrap() error {
	envRoot := os.Getenv(_c.EnvVarMilpaRoot)
	pathMap := NewPathBuilder()

	if envRoot != "" {
		MilpaRoot = envRoot
	} else {
		logrus.Debugf("%s is not set, using default %s", _c.EnvVarMilpaRoot, envRoot)
	}

	if !IsDir(MilpaRoot, false) {
		return errors.EnvironmentError{Err: fmt.Errorf("%s (%s) is not a directory", _c.EnvVarMilpaRoot, MilpaRoot)}
	}

	if len(MilpaPath) != 0 && MilpaPath[0] != "" {
		if isTrueIsh(os.Getenv(_c.EnvVarMilpaPathParsed)) {
			logrus.Debugf("%s already parsed upstream. %d items found", _c.EnvVarMilpaPath, len(MilpaPath))
			return nil
		}

		logrus.Debugf("%s is has %d items, parsing", _c.EnvVarMilpaPath, len(MilpaPath))
		for idx, p := range MilpaPath {
			if p == "" || !IsDir(p, true) {
				logrus.Debugf("Dropping non-directory <%s> from MILPA_PATH", p)
				MilpaPath = append(MilpaPath[:idx], MilpaPath[idx+1:]...)
				continue
			}

			if !strings.HasSuffix(p, _c.RepoRoot) {
				p = filepath.Join(p, _c.RepoRoot)
				logrus.Debugf("Updated path to %s", p)
			}
			pathMap.Add(0, p)
		}
	}

	rootRepo := filepath.Join(MilpaRoot, _c.RepoRoot)
	if !IsDir(rootRepo, false) {
		return errors.EnvironmentError{Err: fmt.Errorf("milpa's built-in repo at %s is not a directory", rootRepo)}
	}

	pathMap.Add(1, rootRepo)
	if pwd, err := os.Getwd(); err == nil {
		pwdRepo := filepath.Join(pwd, _c.RepoRoot)
		if IsDir(pwdRepo, false) {
			logrus.Debugf("Adding pwd repo %s", pwdRepo)
			pathMap.Add(2, pwdRepo)
		}
	}

	pathMap.AddLookup(_c.EnvVarLookupGitDisabled, lookupGitRepo)
	pathMap.AddLookup(_c.EnvVarLookupUserReposDisabled, lookupUserRepos)
	pathMap.AddLookup(_c.EnvVarLookupGlobalReposDisabled, lookupGlobalRepos)

	MilpaPath = pathMap.Ordered()

	return nil
}

func lookupGitRepo() []string {
	logrus.Debugf("looking for a git repo")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	err := cmd.Run()

	if ctx.Err() == nil && err == nil {
		repoRoot := strings.TrimSuffix(stdout.String(), "\n")
		gitRepo := filepath.Join(repoRoot, _c.RepoRoot)
		if IsDir(gitRepo, false) {
			logrus.Debugf("Found repo from git: %s", gitRepo)
			return []string{gitRepo}
		}
	}
	return []string{}
}

func lookupUserRepos() []string {
	logrus.Debugf("looking for user repos")
	found := []string{}
	home := os.Getenv("XDG_DATA_HOME")

	if home == "" {
		home = os.Getenv("HOME")
	}

	if home == "" {
		logrus.Debugf("Ignoring user repo lookup, neither XDG_DATA_HOME nor HOME were found in the environment")
		return found
	}

	userRepos := filepath.Join(home, ".local", "share", "milpa", "repos")
	if files, err := os.ReadDir(userRepos); err == nil {
		for _, file := range files {
			userRepo := filepath.Join(userRepos, file.Name())
			if IsDir(userRepo, true) {
				logrus.Debugf("Found user repo: %s", userRepo)
				found = append(found, userRepo)
			}
		}
	} else {
		logrus.Warnf("User repo directory not found: %s", userRepos)
	}

	return found
}

func lookupGlobalRepos() []string {
	logrus.Debugf("looking for global repos")
	found := []string{}
	globalRepos := filepath.Join(MilpaRoot, "repos")
	if files, err := os.ReadDir(globalRepos); err == nil {
		for _, file := range files {
			globalRepo := filepath.Join(globalRepos, file.Name())
			if IsDir(globalRepo, true) {
				logrus.Debugf("Found global repo: %s", globalRepo)
				found = append(found, globalRepo)
			}
		}
	}

	return found
}
