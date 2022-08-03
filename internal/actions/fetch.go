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
package actions

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NormalizeRepoURI(src string) (uri *url.URL, scheme string, err error) {
	fqurl, err := getter.Detect(src, os.Getenv("PWD"), getter.Detectors)
	if err != nil {
		err = fmt.Errorf("could not detect source for uri %s: %w", src, err)
		return
	}
	logrus.Debugf("Detected uri: %s", fqurl)

	uri, err = url.Parse(fqurl)
	if err != nil {
		return nil, "", fmt.Errorf("could not parse uri %s: %w", fqurl, err)
	}

	scheme = uri.Scheme
	if scheme == "file" || uri.Opaque == "" || uri.Opaque[0] != ':' {
		return
	}

	// remove git:: from git::ssh://git@github.com/...
	logrus.Debugf("Unwrapping uri: %s", uri.Opaque[1:])
	uri, err = url.Parse(uri.Opaque[1:])
	if err != nil {
		err = fmt.Errorf("could not parse unwrapped URI %s: %w", uri.String(), err)
		return
	}

	return
}

func RepoFolderName(uri *url.URL) string {
	path := strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(uri.Path, "/"), "/.milpa"), ".git", "")
	return strings.ReplaceAll(
		strings.Join([]string{
			strings.ReplaceAll(uri.Host, ".", "-"),
			strings.ReplaceAll(path, "/", "-"),
		}, "-"),
		"--",
		"-",
	)
}

var fetchRemoteRepo = &cobra.Command{
	Use:               "__fetch [dst] [src]",
	Short:             "Fetches repos using go-getter",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		src := args[0]
		dst := args[1]

		uri, scheme, err := NormalizeRepoURI(src)
		if err != nil {
			logrus.Fatal(err)
		}

		if scheme == "file" {
			logrus.Fatal("Refusing to copy local folder")
			// fmt.Print(dst + "/" + folder)
			// os.Exit(3)
			// return nil
		}

		folder := RepoFolderName(uri)
		logrus.Debugf("Downloading %s to %s using %s", src, dst+"/"+folder, scheme)

		err = getter.Get(dst+"/"+folder, src)
		if err != nil {
			return err
		}
		logrus.Debug("Download complete")
		fmt.Print(dst + "/" + folder)
		return nil
	},
}
