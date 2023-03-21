// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package actions_test

import (
	"testing"

	"github.com/unrob/milpa/internal/actions"
)

func TestFetchDetectors(t *testing.T) {
	tests := []struct {
		src    string
		dst    string
		folder string
	}{
		{
			"git@github.com:unRob/milpa.git",
			"ssh://git@github.com/unRob/milpa.git",
			"github-com-unRob-milpa",
		},
		{
			"github.com/unRob/milpa",
			"https://github.com/unRob/milpa.git",
			"github-com-unRob-milpa",
		},
		{
			"https://github.com/unRob/milpa.tgz",
			"https://github.com/unRob/milpa.tgz",
			"github-com-unRob-milpa",
		},
	}

	for _, data := range tests {
		t.Run(data.src, func(t *testing.T) {
			uri, _, err := actions.NormalizeRepoURI(data.src)
			if err != nil {
				t.Fatal(err)
			}

			if uri.String() != data.dst {
				t.Fatalf("Wanted uri %s got %s", data.dst, uri)
			}

			folder := actions.RepoFolderName(uri)
			if folder != data.folder {
				t.Fatalf("Wanted folder %s got %s", data.folder, folder)
			}
		})
	}
}
