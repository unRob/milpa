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
package internal

var UsageTPL = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`

var HelpTemplate = `
{{ printf "milpa %s" .Annotations["fullName"] | bold }} - {{ .Short }}

{{ "Usage:" | bold }}

  {{ .Use }}{{if .HasAvailableSubCommands}} subcommand{{end}}

{{- if .HasAvailableSubCommands-}}
{{ "Subcommands:" | bold }}

{{- range .Commands -}}
{{- if (or .IsAvailableCommand (eq .Name "help")) -}}
{{rpad .Name .NamePadding }} {{.Short}}
{{ end -}}
{{- end -}}
{{- end -}}

{{- if hasArguments -}}
{{ "Arguments:" | bold }}

{{ range arguments -}}
  {{- if not .Required}}[{{ end -}}
  ${{- .Name | ToUpper -}}
  {{- if .Variadic}}...{{ end -}}
  {{- if not .Required }}]{{ end }} - {{ .Description }}
{{ end -}}
{{- end -}}

{{- if .HasAvailableLocalFlags}}
{{ "Options:" | bold }}

{{ range $name, $opt := options -}}
- {{printf "--%s" $name | bold}} ({{$opt.Type}}): {{$opt.Description}}.{{ if $opt.Default }} Default: {{$opt.Default}}.{{end}}
{{ end -}}
{{- end -}}
{{- end -}}

{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end -}}

## milpa {{.FullName }}
{{- if .Options }} [options]
{{- end }}
{{- if .Arguments }} {{ range .Arguments -}}
  {{- if not .Required}}[{{ end -}}
  ${{- .Name | ToUpper -}}
  {{- if .Variadic}}...{{ end -}}
  {{- if not .Required }}]{{ end }} {{ end -}}
{{- end}}

{{.Summary}}

{{- if .Arguments }}

## Arguments

{{ range .Arguments -}}
- {{ .Name | ToUpper | AsCode }}{{ if .Variadic }}...{{end}}{{ if .Required }} [REQUIRED]{{ end }}: {{.Description}}.{{ if .Validates }} Valid options{{ if .Set.From.SubCommand }} come from the output of running {{printf "milpa %s" .Set.From.SubCommand | AsCode}}{{end}}{{ if .Set.Values }} are: .Set.Resolved {{end}}.{{end}}
{{ end -}}
{{- else }}
{{ end -}}
{{- if .Options }}
## Options

{{ range $name, $opt := .Options -}}
- {{printf "--%s" $name | AsCode}} ({{$opt.Type}}): {{$opt.Description}}.{{ if $opt.Default }} Default: {{$opt.Default}}.{{end}}
{{ end -}}
{{- end }}
## Description

{{.Description}}`
