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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	_c "github.com/unrob/milpa/internal/constants"
	"github.com/unrob/milpa/internal/errors"
)

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	logrus.Warnf("Discarding path <%s> since its not a directory", path)
	return false
}

var MilpaRoot = "/usr/local/lib/milpa"

type pathLayer map[string]bool

func (pl pathLayer) add(path string) {
	if _, inMap := pl[path]; !inMap {
		pl[path] = true
	}
}

type pathBuilder struct {
	layers map[int]*pathLayer
	unique map[string]bool
	mutex  sync.Mutex
}

func (pb *pathBuilder) add(layerID int, path string, verifyDirectory bool) {
	if pb.unique == nil {
		pb.unique = map[string]bool{}
	}
	if _, exists := pb.unique[path]; exists {
		return
	}
	pb.unique[path] = true

	if verifyDirectory && !isDir(path) {
		logrus.Warnf("Discarding path <%s> since its not a directory", path)
		return
	}

	pb.mutex.Lock()

	if _, exists := pb.layers[layerID]; !exists {
		pb.layers[layerID] = &pathLayer{}
	}

	pb.layers[layerID].add(path)
	pb.mutex.Unlock()
}

func (pb *pathBuilder) Ordered() []string {
	res := []string{}
	keys := []int{}
	for key := range pb.layers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, key := range keys {
		layer := pb.layers[key]
		group := []string{}
		for path := range *layer {
			group = append(group, path)
		}
		sort.Strings(group)
		res = append(res, group...)
	}

	return res
}

func Bootstrap() error {
	envRoot := os.Getenv(_c.EnvVarMilpaRoot)
	pathMap := &pathBuilder{layers: map[int]*pathLayer{}}

	if envRoot != "" {
		MilpaRoot = envRoot
	} else {
		logrus.Debugf("%s is not set, using default %s", _c.EnvVarMilpaRoot, envRoot)
	}

	if !isDir(MilpaRoot) {
		return errors.ConfigError{Err: fmt.Errorf("%s is not a directory", _c.EnvVarMilpaRoot)}
	}

	if len(MilpaPath) != 0 {
		logrus.Debugf("%s is has %d items, parsing", _c.EnvVarMilpaPath, len(MilpaPath))
		for idx, p := range MilpaPath {
			if p == "" || !isDir(p) {
				logrus.Debugf("Dropping item <%s>", p)
				MilpaPath = append(MilpaPath[:idx], MilpaPath[idx+1:]...)
				continue
			}

			pathMap.add(0, p, false)
			if !strings.HasSuffix(p, _c.RepoRoot) {
				p = filepath.Join(p, _c.RepoRoot)
			}
			logrus.Debugf("Updated path to %s", p)
		}
	}

	pathMap.add(1, filepath.Join(MilpaRoot, _c.RepoRoot), true)
	if pwd, err := os.Getwd(); err != nil {
		pwdRepo := filepath.Join(pwd, _c.RepoRoot)
		logrus.Debugf("Adding pwd repo %s", pwdRepo)
		pathMap.add(2, pwdRepo, true)
	}

	lookups := []func(pm *pathBuilder, layer int){}
	if os.Getenv("MILPA_DISABLE_GIT") != "true" {
		lookups = append(lookups, lookupGitRepo)
		// lookupGitRepo(pathMap)
	}
	lookups = append(lookups, lookupUserRepos, lookupGlobalRepos)

	var wg sync.WaitGroup
	for idx, lookup := range lookups {
		wg.Add(1)
		lookup := lookup
		layerID := idx + 10
		go func() {
			defer wg.Done()
			lookup(pathMap, layerID)
		}()
	}

	wg.Wait()
	MilpaPath = pathMap.Ordered()

	return nil
}

func lookupGitRepo(pathMap *pathBuilder, layer int) {
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
		logrus.Debugf("Adding git repo %s", gitRepo)
		pathMap.add(layer, gitRepo, true)
	} // else {
	// logrus.Debugf(ctx.Err())
	// logrus.Error(err)
	// }
}

func lookupUserRepos(pathMap *pathBuilder, layer int) {
	home := os.Getenv("XDG_DATA_HOME")
	if home == "" {
		home = os.Getenv("HOME")
	}

	if home == "" {
		return
	}

	userRepos := filepath.Join(home, ".local", "share", "milpa", "repos")
	if files, err := ioutil.ReadDir(userRepos); err == nil {
		for _, file := range files {
			userRepo := filepath.Join(userRepos, file.Name())
			logrus.Debugf("Adding user repo %s", userRepo)
			pathMap.add(layer, userRepo, false)
		}
	}
}

func lookupGlobalRepos(pathMap *pathBuilder, layer int) {
	globalRepos := filepath.Join(MilpaRoot, "repos")
	if files, err := ioutil.ReadDir(globalRepos); err == nil {
		for _, file := range files {
			globalRepo := filepath.Join(globalRepos, file.Name())
			logrus.Debugf("Adding global repo %s", globalRepo)
			pathMap.add(layer, globalRepo, false)
		}
	}
}