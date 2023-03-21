// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package lookup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"git.rob.mx/nidito/chinampa"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/unrob/milpa/internal/bootstrap"
	"github.com/unrob/milpa/internal/command"
	_c "github.com/unrob/milpa/internal/constants"
)

var log = logger.Sub("lookup")

var DefaultFS = os.DirFS("/")

func Scripts(query []string) (results map[string]string, err error) {
	if len(bootstrap.MilpaPath) == 0 {
		err = fmt.Errorf("no %s set on the environment", _c.EnvVarMilpaPath)
		return
	}

	log.Debugf("looking for scripts in %s=%s", _c.EnvVarMilpaPath, strings.Join(bootstrap.MilpaPath, ":"))
	results = map[string]string{}
	for _, path := range bootstrap.MilpaPath {
		queryBase := strings.Join(append([]string{strings.TrimPrefix(path, "/"), _c.RepoCommandFolderName}, query...), "/")
		matches, err := doublestar.Glob(DefaultFS, fmt.Sprintf("%s/*", queryBase), doublestar.WithFilesOnly())

		if err != nil {
			log.Debugf("errored while globbing")
			continue
		}

		log.Debugf("found %d potential matches in %s", len(matches), path)
		for _, match := range matches {
			extension := filepath.Ext(match)
			if extension != "" && extension != ".sh" {
				if extension != ".yaml" {
					log.Debugf("ignoring /%s, unknown extension", match)
				}
				continue
			}

			results["/"+match] = path
		}
	}

	return results, err
}

func AllSubCommands(returnOnError bool) error {
	files, err := Scripts([]string{"**"})
	if err != nil {
		return err
	}

	log.Debugf("Found %d files", len(files))

	// make sure we always sort commands by path before initializing
	// this helps with "index" commands, i.e. commands named like an existing folder
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, path := range keys {
		repo := files[path]
		cmd, specErr := command.New(path, repo)
		if specErr != nil {
			if returnOnError {
				return specErr
			}
		}

		log.Debugf("Initialized %s", cmd.FullName())
		chinampa.Register(cmd)
	}

	return err
}

func AllDocs() ([]string, error) {
	results := []string{}
	if err := bootstrap.CheckMilpaPathSet(); err != nil {
		return results, err
	}

	log.WithField("kind", "docs").Debugf("looking for all docs in %s", bootstrap.MilpaPath)

	for _, path := range bootstrap.MilpaPath {
		q := path + "/" + _c.RepoDocsFolderName + "/**/*.md"

		log.WithField("kind", "docs").Debugf("looking for all docs matching %s", q)
		basepath, pattern := doublestar.SplitPattern(q)
		fsys := os.DirFS(basepath)
		docs, err := doublestar.Glob(fsys, pattern, doublestar.WithFilesOnly())
		if err != nil {
			log.WithField("kind", "docs").Debugf("errored looking for all docs matching %s: %s", q, err)
			return results, err
		}

		log.WithField("kind", "docs").Debugf("found %d docs matching %s", len(docs), q)

		for _, doc := range docs {
			if strings.Contains(doc, _c.RepoDocsTemplateFolderName) {
				continue
			}
			results = append(results, basepath+"/"+doc)
		}
	}

	return results, nil
}

func Docs(query []string, needle string, returnPaths bool) ([]string, error) {
	results := []string{}
	found := map[string]bool{}
	if err := bootstrap.CheckMilpaPathSet(); err != nil {
		return results, err
	}

	log.WithField("kind", "docs").Debugf("looking for docs in %s with", bootstrap.MilpaPath)
	queryString := ""
	if len(query) > 0 {
		queryString = strings.Join(query, "/")
	}

	for _, path := range bootstrap.MilpaPath {
		qbase := path + "/" + _c.RepoDocsFolderName
		if len(query) > 0 {
			qbase += "/" + queryString
		}
		q := qbase + "/*"
		if returnPaths {
			q = qbase + "/*.md"
		}
		log.WithField("kind", "docs").Debugf("looking for docs matching %s", q)
		docs, err := filepath.Glob(q)
		if err != nil {
			log.WithField("kind", "docs").Debugf("failed looking for docs matching %s", q)
			return results, err
		}

		log.WithField("kind", "docs").Debugf("Found %d docs matching %s", len(docs), q)
		for _, doc := range docs {
			fname := filepath.Base(doc)
			extensionParts := strings.Split(fname, ".")
			ext := ""
			if len(extensionParts) > 1 {
				ext = extensionParts[len(extensionParts)-1]
			}

			if strings.Contains(doc, "/"+_c.RepoDocsTemplateFolderName) || (ext != "" && ext != "md") {
				log.WithField("kind", "docs").Debugf("Ignoring non-doc file: %s, ext: %s, md: %v", doc, ext, (ext != "" && ext != ".md"))
				continue
			}
			name := strings.TrimSuffix(fname, ".md")
			if _, ok := found[name]; (needle == "" || strings.HasPrefix(name, needle)) && !ok {
				if returnPaths {
					results = append(results, doc)
				} else {
					results = append(results, name)
				}
				found[name] = true
			} else {
				log.WithField("kind", "docs").Debugf("ignoring %s => %s", name, doc)
			}
		}
	}

	log.WithField("kind", "docs").Debugf("returning %d docs", len(results))
	return results, nil
}
