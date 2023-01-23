// SPDX-License-Identifier: Apache-2.0
// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
package constants

import (
	"regexp"

	// Embed requires an import so the compiler knows what's up. Golint requires a comment. Gotta please em both.
	_ "embed"

	"git.rob.mx/nidito/chinampa/pkg/env"
)

const Milpa = "milpa"
const HelpCommandName = "help"

func init() {
	env.HelpUnstyled = "MILPA_PLAIN_HELP"
	env.HelpStyle = "MILPA_HELP_STYLE"
	env.Verbose = "MILPA_VERBOSE"
	env.Silent = "MILPA_SILENT"
	env.ValidationDisabled = "MILPA_SKIP_VALIDATION"
}

// Environment Variables.
const EnvVarMilpaPath = "MILPA_PATH"
const EnvVarMilpaPathParsed = "MILPA_PATH_PARSED"
const EnvVarMilpaRoot = "MILPA_ROOT"
const EnvVarCompaOut = "COMPA_OUT"
const EnvVarDebug = "DEBUG"
const EnvVarLookupGitDisabled = "MILPA_DISABLE_GIT"
const EnvVarLookupUserReposDisabled = "MILPA_DISABLE_USER_REPOS"
const EnvVarLookupGlobalReposDisabled = "MILPA_DISABLE_GLOBAL_REPOS"

// Folder structure.
const RepoRoot = ".milpa"
const RepoCommandFolderName = "commands"
const RepoCommands = ".milpa/commands"
const RepoDocsFolderName = "docs"
const RepoDocsTemplateFolderName = ".template"
const RepoDocs = ".milpa/docs"

// Output variable prefixes.
const OutputPrefixArg = "MILPA_ARG_"
const OutputPrefixOpt = "MILPA_OPT_"
const OutputCommandName = "MILPA_COMMAND_NAME"
const OutputCommandKind = "MILPA_COMMAND_KIND"
const OutputCommandRepo = "MILPA_COMMAND_REPO"
const OutputCommandPath = "MILPA_COMMAND_PATH"

var OutputPrefixPattern = regexp.MustCompile(`\$\{?[#!]?MILPA_((OPT|ARG)_([0-9a-zA-Z_]+))`)
