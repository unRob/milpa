// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package docs

import (
	"fmt"
	"os"
	"strings"

	"git.rob.mx/nidito/chinampa/pkg/errors"
	"git.rob.mx/nidito/chinampa/pkg/logger"
	"github.com/unrob/milpa/internal/bootstrap"
)

var log = logger.Sub("documentation")

func FromQuery(query []string) ([]byte, error) {
	if err := bootstrap.CheckMilpaPathSet(); err != nil {
		return []byte{}, err
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("requesting docs help")
	}

	queryString := strings.Join(query, "/")

	for _, path := range bootstrap.MilpaPath {
		candidate := path + "/docs/" + queryString
		log.Debugf("looking for doc named %s", candidate)
		_, err := os.Lstat(candidate + ".md") // nolint: gosec
		if err == nil {
			log.Debugf("found doc for %s", candidate)
			return os.ReadFile(candidate + ".md") // nolint: gosec
		}

		if _, err := os.Lstat(candidate + "/index.md"); err == nil { // nolint: gosec
			log.Debugf("found index doc for %s", candidate)
			return os.ReadFile(candidate + "/index.md") // nolint: gosec
		}

		if _, err := os.Stat(candidate); err == nil { // nolint: gosec
			break
		}
	}

	missingPath := strings.Join(query, "/")
	return nil, errors.BadArguments{Msg: fmt.Sprintf("Missing topic named <%s.md> or <%s/index.md> in any of %s", missingPath, missingPath, strings.Join(bootstrap.MilpaPath, ":"))}
}
