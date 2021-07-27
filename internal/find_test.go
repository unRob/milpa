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
package internal_test

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/sirupsen/logrus"
	. "github.com/unrob/milpa/internal"
)

var fsBase = "usr/local/milpa"
var allCommands = map[string]*fstest.MapFile{
	"shell-script.sh": {
		Data: []byte(`#!/usr/bin env bash
echo "hello"`),
		Mode: 0644,
	},
	"shell-script.yaml": {
		Data: []byte(`{"description": "a lot of stuff", "summary": "shell scripts stuff"}`),
		Mode: 0644,
	},
	"nested/shell-script.sh": {
		Data: []byte(`#!/usr/bin env bash
echo "hello"`),
		Mode: 0644,
	},
	"nested/shell-script.yaml": {
		Data: []byte(`{"description": "a lot of stuff", "summary": "shell scripts stuff"}`),
		Mode: 0644,
	},
	"executable": {
		Data: []byte(`some bytes`),
		Mode: 0777,
	},
	"executable.yaml": {
		Data: []byte(`{"description": "a lot of stuff", "summary": "execs stuff"}`),
	},
	"nested/executable": {
		Data: []byte(`some bytes`),
		Mode: 0777,
	},
	"nested/executable.yaml": {
		Data: []byte(`{"description": "a lot of stuff", "summary": "execs stuff"}`),
	},
	"misleading": {
		Data: []byte(`some bytes`),
		Mode: 0777,
	},
	"misleading.sh": {
		Data: []byte(`some bytes`),
		Mode: 0644,
	},
	"misleading.yaml": {
		Data: []byte(`{"description": "a lot of stuff", "summary": "execs stuff"}`),
	},
	"missing-spec": {
		Data: []byte(`some bytes`),
		Mode: 0777,
	},
	"empty-dir": {
		Mode: fs.ModeDir,
	},
}

func setupFS(filenames []string, pool map[string]*fstest.MapFile) *fstest.MapFS {
	fs := fstest.MapFS{
		fsBase + "/.milpa":                 {Mode: fs.ModeDir},
		fsBase + "/.milpa/commands":        {Mode: fs.ModeDir},
		fsBase + "/.milpa/commands/nested": {Mode: fs.ModeDir},
	}

	for _, name := range filenames {
		fs[fsBase+"/.milpa/commands/"+name] = pool[name]
	}

	MilpaPath = []string{fsBase + "/.milpa"}
	DefaultFS = &fs
	return &fs
}

func TestFindScripts(t *testing.T) {
	selected := []string{
		"shell-script.sh",
		"shell-script.yaml",
		"executable",
		"executable.yaml",
		"empty-dir",
		"missing-spec",
		"misleading",
		"misleading.sh",
		"misleading.yaml",
		"nested/shell-script.sh",
		"nested/shell-script.yaml",
		"nested/executable",
		"nested/executable.yaml",
	}
	mfs := setupFS(selected, allCommands)
	logrus.SetLevel(logrus.DebugLevel)
	files, err := FindScripts([]string{"**"})
	if err != nil {
		t.Fatalf("Could not find scripts: %v", err)
	}

	expectedPaths := []string{
		"shell-script.sh",
		"executable",
		"missing-spec",
		"misleading",
		"misleading.sh",
		"nested/shell-script.sh",
		"nested/executable",
	}

	if len(files) != len(expectedPaths) {
		t.Errorf("Found incorrect amount of scripts: %d vs %d; %v", len(files), len(expectedPaths), files)
	}

	expected := map[string]os.FileInfo{}
	for _, p := range expectedPaths {
		fp := fsBase + "/.milpa/commands/" + p
		i, err := mfs.Stat(fp)
		if err != nil {
			t.Fatalf("failed setting expected results (%s): %v", p, err)
		}

		expected["/"+fp] = i
	}

	for efile, einfo := range expected {
		data, exists := files[efile]
		if !exists {
			t.Errorf("missing file %s", efile)
			continue
		}
		delete(files, efile)

		if data.Info.Mode() != einfo.Mode() {
			t.Errorf("Unexpected mode. Expected: %v, got: %v", einfo.Mode(), data.Info.Mode())
		}
	}
}
