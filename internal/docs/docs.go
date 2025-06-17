// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/unrob/milpa/internal/repo"
)

var log = logger.Sub("documentation")

func FromQuery(query []string) ([]byte, error) {
	if err := repo.CheckPathSet(); err != nil {
		return []byte{}, err
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("requesting docs help")
	}

	queryString := strings.Join(query, "/")

	for _, path := range repo.Path {
		candidate := path + "/docs/" + queryString
		log.Debugf("looking for doc named %s", candidate)
		_, err := os.Lstat(candidate + ".md")
		if err == nil {
			log.Debugf("found doc for %s", candidate)
			return os.ReadFile(candidate + ".md")
		}

		if _, err := os.Lstat(candidate + "/index.md"); err == nil {
			log.Debugf("found index doc for %s", candidate)
			return os.ReadFile(candidate + "/index.md")
		}

		if _, err := os.Stat(candidate); err == nil {
			break
		}
	}

	missingPath := strings.Join(query, "/")
	return nil, errors.BadArguments{Msg: fmt.Sprintf("Missing topic named <%s.md> or <%s/index.md> in any of %s", missingPath, missingPath, strings.Join(repo.Path, ":"))}
}
