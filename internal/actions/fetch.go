// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/command"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/hashicorp/go-getter"
)

var fetchLog = logger.Sub("itself repo install")

func NormalizeRepoURI(src string) (uri *url.URL, scheme string, err error) {
	fqurl, err := getter.Detect(src, os.Getenv("PWD"), getter.Detectors)
	if err != nil {
		err = fmt.Errorf("could not detect source for uri %s: %w", src, err)
		return
	}
	fetchLog.Debugf("Detected uri: %s", fqurl)

	uri, err = url.Parse(fqurl)
	if err != nil {
		return nil, "", fmt.Errorf("could not parse uri %s: %w", fqurl, err)
	}

	scheme = uri.Scheme
	if scheme == "file" || uri.Opaque == "" || uri.Opaque[0] != ':' {
		return
	}

	// remove git:: from git::ssh://git@github.com/...
	fetchLog.Debugf("Unwrapping uri: %s", uri.Opaque[1:])
	uri, err = url.Parse(uri.Opaque[1:])
	if err != nil {
		err = fmt.Errorf("could not parse unwrapped URI %s: %w", uri.String(), err)
		return
	}

	return
}

var extensions = []string{
	"/",
	".tar.gz",
	".tgz",
	".tar.bz2",
	".tbz2",
	".tar.xz",
	".txz",
	".zip",
	".gz",
	".bz2",
	".xz",
	".git",
	"/.milpa",
}

func removeExtensions(path string) string {
	for _, suffix := range extensions {
		path = strings.TrimSuffix(path, suffix)
	}
	return path
}

func RepoFolderName(uri *url.URL) string {
	return strings.ReplaceAll(
		strings.Join([]string{
			strings.ReplaceAll(uri.Host, ".", "-"),
			strings.ReplaceAll(removeExtensions(uri.Path), "/", "-"),
		}, "-"),
		"--",
		"-",
	)
}

var Fetch = &command.Command{
	Path:        []string{"__fetch"},
	Hidden:      true,
	Summary:     "Fetches repos using go-getter",
	Description: `Yep`,
	Arguments: command.Arguments{
		{
			Name:        "source",
			Description: "The source to fetch from",
			Required:    true,
		},
		{
			Name:        "target",
			Description: "The destination to fetch to",
			Required:    true,
		},
	},
	Action: func(cmd *command.Command) error {
		src := cmd.Arguments[0].ToString()
		dst := cmd.Arguments[1].ToString()

		uri, scheme, err := NormalizeRepoURI(src)
		if err != nil {
			fetchLog.Fatal(err)
		}

		if scheme == "file" {
			fetchLog.Fatal("Refusing to copy local folder")
		}

		folder := RepoFolderName(uri)
		fetchLog.Debugf("Downloading %s to %s using %s", src, dst+"/"+folder, scheme)

		err = getter.Get(dst+"/"+folder, src)
		if err != nil {
			return err
		}
		fetchLog.Debug("Download complete")
		fmt.Print(dst + "/" + folder)
		return nil
	},
}
