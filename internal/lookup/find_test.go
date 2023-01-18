// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package lookup_test

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/sirupsen/logrus"
	"github.com/unrob/milpa/internal/bootstrap"
	. "github.com/unrob/milpa/internal/lookup"
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

var noDocs = map[string]*fstest.MapFile{}

func setupFS(filenames []string, pool map[string]*fstest.MapFile, docs map[string]*fstest.MapFile) *fstest.MapFS {
	fs := fstest.MapFS{
		fsBase + "/.milpa":                 {Mode: fs.ModeDir},
		fsBase + "/.milpa/commands":        {Mode: fs.ModeDir},
		fsBase + "/.milpa/commands/nested": {Mode: fs.ModeDir},
	}

	for _, name := range filenames {
		fs[fsBase+"/.milpa/commands/"+name] = pool[name]
	}

	for path, f := range docs {
		fs[fsBase+"/.milpa/docs/"+path] = f
	}

	bootstrap.MilpaPath = []string{fsBase + "/.milpa"}
	DefaultFS = &fs
	return &fs
}

func TestScripts(t *testing.T) {
	t.Run("errors without milpa_path set", func(t *testing.T) {
		mp := bootstrap.MilpaPath
		defer func() { bootstrap.MilpaPath = mp }()
		bootstrap.MilpaPath = []string{}
		if _, err := Scripts([]string{"**"}); err == nil {
			t.Fatalf("did not error as expected")
		}
	})

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
	mfs := setupFS(selected, allCommands, noDocs)
	logrus.SetLevel(logrus.DebugLevel)
	files, err := Scripts([]string{"**"})
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

func TestDocsFind(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	root := fromProjectRoot()
	bootstrap.MilpaPath = []string{root + "/.milpa"}

	t.Run("top-level", func(t *testing.T) {
		topics, err := Docs([]string{}, "", false)
		if err != nil {
			t.Fatalf("did not find docs: %s", err)
		}

		expected := []string{"milpa"}
		if len(topics) != len(expected) || fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", topics) {
			t.Fatalf("Did not find expected docs:\nwanted: %s\ngot: %s", expected, topics)
		}
	})

	t.Run("sub-level", func(t *testing.T) {
		topics, err := Docs([]string{"milpa"}, "", false)
		if err != nil {
			t.Fatalf("did not find docs: %s", err)
		}

		expected := []string{"command", "environment", "index", "internals", "quick-guide", "repo", "support", "use-case", "util"}
		if len(topics) != len(expected) || fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", topics) {
			t.Fatalf("Did not find expected docs:\nwanted: %s\ngot: %s", expected, topics)
		}
	})

	t.Run("sub-level autocomplete", func(t *testing.T) {
		topics, err := Docs([]string{"milpa"}, "env", false)
		if err != nil {
			t.Fatalf("did not find docs: %s", err)
		}

		expected := []string{"environment"}
		if len(topics) != len(expected) || fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", topics) {
			t.Fatalf("Did not find expected docs:\nwanted: %s\ngot: %s", expected, topics)
		}
	})

	t.Run("sub-level autocomplete files", func(t *testing.T) {
		topics, err := Docs([]string{"milpa"}, "env", true)
		if err != nil {
			t.Fatalf("did not find docs: %s", err)
		}

		expected := []string{root + "/.milpa/docs/milpa/environment.md"}
		if len(topics) != len(expected) || fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", topics) {
			t.Fatalf("Did not find expected docs:\nwanted: %s\ngot: %s", expected, topics)
		}
	})
}

func TestDocsFindAll(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	root := fromProjectRoot()
	bootstrap.MilpaPath = []string{root + "/.milpa"}

	paths, err := AllDocs()
	if err != nil {
		t.Fatalf("did not find docs: %s", err)
	}

	expected := []string{
		root + "/.milpa/docs/milpa/environment.md",
		root + "/.milpa/docs/milpa/index.md",
		root + "/.milpa/docs/milpa/internals.md",
		root + "/.milpa/docs/milpa/quick-guide.md",
		root + "/.milpa/docs/milpa/support.md",
		root + "/.milpa/docs/milpa/use-case.md",
		root + "/.milpa/docs/milpa/command/index.md",
		root + "/.milpa/docs/milpa/command/spec.md",
		root + "/.milpa/docs/milpa/repo/docs.md",
		root + "/.milpa/docs/milpa/repo/hooks.md",
		root + "/.milpa/docs/milpa/repo/index.md",
		root + "/.milpa/docs/milpa/util/index.md",
		root + "/.milpa/docs/milpa/util/log.md",
		root + "/.milpa/docs/milpa/util/repo.md",
		root + "/.milpa/docs/milpa/util/shell.md",
		root + "/.milpa/docs/milpa/util/tmp.md",
	}
	if len(paths) != len(expected) || fmt.Sprintf("%s", expected) != fmt.Sprintf("%s", paths) {
		t.Fatalf("Did not find expected docs:\nwanted: %s\ngot: %s", expected, paths)
	}
}
