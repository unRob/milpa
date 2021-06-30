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
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/sirupsen/logrus"
)

var MilpaPath []string = strings.Split(os.Getenv("MILPA_PATH"), ":")
var DefaultFS = os.DirFS("/")

type findFilterfunc func(match string, info os.FileInfo) bool

func FilterNoneFunc(match string, info os.FileInfo) bool {
	return false
}

func FilterDirectoriesFunc(match string, info os.FileInfo) bool {
	return !info.IsDir()
}

type FindResult struct {
}

func FindScripts(query []string, filter findFilterfunc) (results map[string]struct {
	Info os.FileInfo
	Repo string
}, err error) {

	if len(MilpaPath) == 0 {
		err = fmt.Errorf("no MILPA_PATH set on the environment")
		return
	}

	logrus.Debugf("looking for scripts in %s", MilpaPath)
	results = map[string]struct {
		Info os.FileInfo
		Repo string
	}{}
	for _, path := range MilpaPath {
		queryBase := strings.Join(append([]string{strings.TrimPrefix(path, "/"), cmdPath}, query...), "/")
		matches, err := doublestar.Glob(DefaultFS, fmt.Sprintf("%s/*", queryBase))

		if err != nil {
			logrus.Debugf("errored while globbing")
			continue
		}

		logrus.Debugf("found %d potential matches in %s", len(matches), path)
		for _, match := range matches {
			if !(strings.HasSuffix(match, ".sh") || filepath.Ext(match) == "") {
				logrus.Debugf("ignoring %s, unknown extension", match)
				continue
			}

			fileInfo, err := fs.Stat(DefaultFS, match)
			if err != nil {
				logrus.Debugf("ignoring %s, failed to stat: %v", match, err)
				continue
			}

			if filter(match, fileInfo) {
				results["/"+match] = struct {
					Info fs.FileInfo
					Repo string
				}{fileInfo, path}
			} else {
				logrus.Debugf("ignoring %s, filtered", match)
			}
		}
	}

	return
}

func FindAllSubCommands(ignoreParsingErrors bool) (cmds []*Command, err error) {
	files, err := FindScripts([]string{"**"}, FilterDirectoriesFunc)
	if err != nil {
		return cmds, err
	}

	logrus.Debugf("Found %d files", len(files))

	for path, data := range files {
		var cmd *Command
		cmd, err = New(path, data.Repo)
		if err != nil {
			logrus.Debugf("Could not initialize command %s", path)
			return
		}
		logrus.Debugf("Initialized %s", cmd.FullName())

		cmds = append(cmds, cmd)
	}

	return
}
