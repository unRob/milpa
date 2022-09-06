{{- if not .HTMLOutput }}
# {{ if and (not (eq .Spec.Meta.Kind "root")) (not (eq .Command.Name "help")) }}milpa {{ end }}{{ .Spec.FullName }}{{if eq .Command.Name "help"}} help{{end}}
{{- else }}
---
description: {{ .Command.Short }}
---
{{- end }}

{{ .Command.Short }}

## Usage

  ﹅{{ replace .Command.UseLine " [flags]" "" }}{{if .Command.HasAvailableSubCommands}} SUBCOMMAND{{end}}﹅

{{ if .Command.HasAvailableSubCommands -}}
## Subcommands

{{ $hh := .HTMLOutput -}}
{{ range .Command.Commands -}}
{{- if (or .IsAvailableCommand (eq .Name "help")) -}}
- {{ if $hh -}}
[﹅{{ .Name }}﹅]({{.Name}})
{{- else -}}
﹅{{ .Name }}﹅
{{- end }} - {{.Short}}
{{ end }}
{{- end -}}
{{- end -}}

{{- if .Spec.Arguments -}}
## Arguments

{{ range .Spec.Arguments -}}

- ﹅{{ .Name | toUpper }}{{ if .Variadic}}...{{ end }}﹅{{ if .Required }} _required_{{ end }} - {{ .Description }}
{{ end -}}
{{- end -}}


{{ if and (eq .Spec.Meta.Kind "root") (not (eq .Command.Name "help")) }}
## Description

{{ .Spec.Description }}
{{ end -}}
{{- if .Spec.HasAdditionalHelp }}
{{ .Spec.AdditionalHelp .HTMLOutput }}
{{ end -}}


{{- if .Command.HasAvailableLocalFlags}}
## Options

{{ range $name, $opt := .Spec.Options -}}
- ﹅--{{ $name }}﹅ (_{{$opt.Type}}_): {{ trimSuffix $opt.Description "."}}.{{ if $opt.Default }} Default: _{{ $opt.Default }}_.{{ end }}
{{ end -}}
{{- end -}}

{{- if not (eq .Command.Name "milpa") }}
## Description

{{ if not (eq .Command.Long "") }}{{ .Command.Long }}{{ else }}{{ .Spec.Description }}{{end}}
{{ end }}

{{- if .Command.HasAvailableInheritedFlags }}
## Global Options

{{ range $name, $opt := .GlobalOptions -}}
- ﹅--{{ $name }}﹅ (_{{$opt.Type}}_): {{$opt.Description}}.{{ if $opt.Default }} Default: _{{ $opt.Default }}_.{{ end }}
{{ end -}}
{{end}}
