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
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/sirupsen/logrus"
)

var MilpaPath []string = strings.Split(os.Getenv("MILPA_PATH"), ":")

type findResult struct {
	Path string
	Info os.FileInfo
}

type findFilterfunc func(match string, info os.FileInfo) bool

func findScripts(query []string, filter findFilterfunc) (results []*findResult) {
	logrus.Debugf("looking for scripts in %s", MilpaPath)
	for _, path := range MilpaPath {
		queryBase := fmt.Sprintf("%s/.milpa/commands/%s", path, strings.Join(query, "/"))
		matches, err := doublestar.Glob(fmt.Sprintf("%s/*{.sh,}", queryBase))

		if err == nil {
			logrus.Debugf("found %d potential matches in %s", len(matches), path)
			for _, match := range matches {
				if !(strings.HasSuffix(match, ".sh") || filepath.Ext(match) == "") {
					logrus.Debugf("ignoring %s", match)
					continue
				}

				fileInfo, err := os.Stat(match)
				if err != nil {
					logrus.Debugf("ignoring %s, failed to stat: %v", match, err)
					continue
				}

				if filter(match, fileInfo) {
					results = append(results, &findResult{match, fileInfo})
				}
			}
		} else {
			logrus.Debugf("errored while globbing")
		}
	}

	return
}

func FindAllSubCommands() (cmds []*Command, err error) {
	files := findScripts([]string{"**"}, func(_ string, info os.FileInfo) bool {
		return !info.IsDir()
	})

	logrus.Debugf("Found %d files", len(files))

	for _, file := range files {
		pc := strings.SplitN(file.Path, "/.milpa/commands/", 2)
		pkg := pc[0]
		kind := ""
		spec := ""

		if strings.HasSuffix(file.Path, ".sh") {
			kind = "source"
			spec = fmt.Sprintf("%s/.milpa/commands/%s.yaml", pkg, strings.Replace(pc[1], ".sh", "", 1))
		} else {
			kind = "exec"
			spec = fmt.Sprintf("%s.yaml", file.Path)
		}

		var cmd *Command
		cmd, err = New(file.Path, spec, pkg, kind)
		if err != nil {
			logrus.Debugf("Could not initialize command %s", file.Path)
			return
		}
		logrus.Debugf("Initialized %s", cmd.FullName())

		cmds = append(cmds, cmd)
	}

	return
}
