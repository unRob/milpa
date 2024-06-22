// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package bootstrap

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"git.rob.mx/nidito/chinampa/pkg/logger"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
	"github.com/unrob/milpa/internal/repo"
	"github.com/unrob/milpa/internal/util"
)

var log = logger.Sub("bootstrap")

func Run() error {
	envRoot := os.Getenv(_c.EnvVarMilpaRoot)

	if envRoot != "" {
		repo.Root = envRoot
	} else {
		log.Debugf("%s is not set, using default %s", _c.EnvVarMilpaRoot, envRoot)
	}

	if !IsDir(repo.Root, false) {
		return errors.EnvironmentError{Err: fmt.Errorf("%s (%s) is not a directory", _c.EnvVarMilpaRoot, repo.Root)}
	}

	pathMap := NewPathBuilder()
	if len(repo.Path) != 0 && repo.Path[0] != "" {
		if util.IsTrueIsh(os.Getenv(_c.EnvVarMilpaPathParsed)) {
			log.Debugf("%s already parsed upstream. %d items found", _c.EnvVarMilpaPath, len(repo.Path))
			return nil
		}

		log.Debugf("%s is has %d items, parsing", _c.EnvVarMilpaPath, len(repo.Path))
		for idx, p := range repo.Path {
			if p == "" || !IsDir(p, true) {
				log.Debugf("Dropping non-directory <%s> from MILPA_PATH", p)
				repo.Path = append(repo.Path[:idx], repo.Path[idx+1:]...)
				continue
			}

			if !strings.HasSuffix(p, _c.RepoRoot) {
				p = filepath.Join(p, _c.RepoRoot)
				log.Debugf("Updated path to %s", p)
			}
			pathMap.Add(0, p)
		}
	}

	rootRepo := filepath.Join(repo.Root, _c.RepoRoot)
	if !IsDir(rootRepo, false) {
		return errors.EnvironmentError{Err: fmt.Errorf("milpa's built-in repo at %s is not a directory", rootRepo)}
	}

	pathMap.Add(1, rootRepo)
	if pwd, err := os.Getwd(); err == nil {
		pwdRepo := filepath.Join(pwd, _c.RepoRoot)
		if IsDir(pwdRepo, false) {
			log.Debugf("Adding pwd repo %s", pwdRepo)
			pathMap.Add(2, pwdRepo)
		}
	}

	pathMap.AddLookup(_c.EnvVarLookupGitDisabled, lookupGitRepo)
	pathMap.AddLookup(_c.EnvVarLookupUserReposDisabled, lookupUserRepos)
	pathMap.AddLookup(_c.EnvVarLookupGlobalReposDisabled, lookupGlobalRepos)

	repo.Path = pathMap.Ordered()
	os.Setenv(_c.EnvVarMilpaPath, strings.Join(repo.Path, ":"))

	return nil
}

func lookupGitRepo() []string {
	log.Debugf("looking for a git repo")
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
			log.Debugf("Found repo from git: %s", gitRepo)
			return []string{gitRepo}
		}
	}
	return []string{}
}

func lookupUserRepos() []string {
	log.Debugf("looking for user repos")
	found := []string{}

	var userRepos string
	if home := os.Getenv("XDG_DATA_HOME"); home == "" {
		home = os.Getenv("HOME")
		userRepos = filepath.Join(home, ".local", "share", "milpa", "repos")
		if home == "" {
			log.Debugf("Ignoring user repo lookup, neither XDG_DATA_HOME nor HOME were found in the environment")
			return found
		}
	} else {
		userRepos = filepath.Join(home, "milpa", "repos")
	}

	if files, err := os.ReadDir(userRepos); err == nil {
		for _, file := range files {
			userRepo := filepath.Join(userRepos, file.Name())
			if IsDir(userRepo, true) {
				log.Debugf("Found user repo: %s", userRepo)
				found = append(found, userRepo)
			}
		}
	} else {
		log.Warnf("User repo directory not found: %s", userRepos)
	}

	return found
}

func lookupGlobalRepos() []string {
	log.Debugf("looking for global repos")
	found := []string{}
	globalRepos := filepath.Join(repo.Root, "repos")
	if files, err := os.ReadDir(globalRepos); err == nil {
		for _, file := range files {
			globalRepo := filepath.Join(globalRepos, file.Name())
			if IsDir(globalRepo, true) {
				log.Debugf("Found global repo: %s", globalRepo)
				found = append(found, globalRepo)
			}
		}
	}

	return found
}
