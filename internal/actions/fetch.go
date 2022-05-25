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

var fetchCommand *cobra.Command = &cobra.Command{
	Use:               "__fetch [dst] [src]",
	Short:             "Fetches stuff using go-getter",
	Hidden:            true,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	Args:              cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		src := args[0]
		dst := args[1]

		fqurl, err := getter.Detect(src, os.Getenv("PWD"), getter.Detectors)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debugf("Detected uri: %s", fqurl)

		uri, err := url.Parse(fqurl)
		if err != nil {
			logrus.Fatal(err)
		}
		scheme := uri.Scheme

		if scheme == "file" {
			logrus.Fatal("Refusing to copy local folder")
		}

		if uri.Opaque != "" && uri.Opaque[0] == ':' {
			logrus.Debugf("Unwrapping uri: %s", uri.Opaque[1:])
			uri2, err := url.Parse(uri.Opaque[1:])
			if err != nil {
				logrus.Fatal(err)
			}

			uri = uri2
		}

		folder := strings.ReplaceAll(uri.Host, ".", "-")
		folder += "-"
		folder += strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(uri.Path, "/"), "/.milpa"), "/", "-")
		folder = strings.ReplaceAll(folder, "--", "-")
		folder = strings.ReplaceAll(folder, ".git", "")

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
