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
package constants

import (
	"regexp"
	"strings"
	"text/template"

	// Embed requires an import so the compiler knows what's up. Golint requires a comment. Gotta please em both.
	_ "embed"
)

const Milpa = "milpa"
const HelpCommandName = "help"

// Environment Variables.
const EnvVarHelpUnstyled = "MILPA_PLAIN_HELP"
const EnvVarHelpStyle = "MILPA_HELP_STYLE"
const EnvVarMilpaRoot = "MILPA_ROOT"
const EnvVarMilpaPath = "MILPA_PATH"
const EnvVarMilpaPathParsed = "MILPA_PATH_PARSED"
const EnvVarMilpaVerbose = "MILPA_VERBOSE"
const EnvVarMilpaSilent = "MILPA_SILENT"
const EnvVarMilpaUnstyled = "NO_COLOR"
const EnvVarValidationDisabled = "MILPA_SKIP_VALIDATION"
const EnvVarCompaOut = "COMPA_OUT"
const EnvVarDebug = "DEBUG"
const EnvVarLookupGitDisabled = "MILPA_DISABLE_GIT"
const EnvVarLookupUserReposDisabled = "MILPA_DISABLE_USER_REPOS"
const EnvVarLookupGlobalReposDisabled = "MILPA_DISABLE_GLOBAL_REPOS"

// EnvFlagNames are flags also available as environment variables.
var EnvFlagNames = map[string]string{
	"no-color":        EnvVarMilpaUnstyled,
	"silent":          EnvVarMilpaSilent,
	"verbose":         EnvVarMilpaVerbose,
	"skip-validation": EnvVarValidationDisabled,
}

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

// Exit statuses
// see man sysexits || grep "#define EX" /usr/include/sysexits.h
// and https://tldp.org/LDP/abs/html/exitcodes.html

// 0 means everything is fine.
const ExitStatusOk = 0

// 42 provides answers to life, the universe and everything; also, renders help.
const ExitStatusRenderHelp = 42

// 64 bad arguments
// EX_USAGE The command was used incorrectly, e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
const ExitStatusUsage = 64

// EX_SOFTWARE An internal software error has been detected. This should be limited to non-operating system related errors as possible.
const ExitStatusProgrammerError = 70

// EX_CONFIG Something was found in an unconfigured or misconfigured state.
const ExitStatusConfigError = 78

// 127 command not found.
const ExitStatusNotFound = 127

// ContextKeyRuntimeIndex is the string key used to store context in a cobra Command.
const ContextKeyRuntimeIndex = "x-milpa-runtime-index"

//go:embed help.md
var helpTemplateText string

// TemplateCommandHelp holds a template for rendering command help.
var TemplateCommandHelp = template.Must(template.New("help").Funcs(template.FuncMap{
	"trim":       strings.TrimSpace,
	"toUpper":    strings.ToUpper,
	"trimSuffix": strings.TrimSuffix,
}).Parse(helpTemplateText))
