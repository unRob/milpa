package command

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/cpuguy83/go-md2man/md2man"
)

var HELP_TEMPLATE = `## milpa {{.FullName }}
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

var MAN_TEMPLATE = `# "milpa {{.FullName}}" 42 "some other text" "and more"
` + HELP_TEMPLATE

func (cmd *Command) Help(format string) ([]byte, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"AsCode": func(thing interface{}) string {
			return fmt.Sprintf("`%s`", thing)
		},
	}

	tpl, err := template.New("help-template").Funcs(funcMap).Parse(HELP_TEMPLATE)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, cmd); err != nil {
		return nil, err
	}

	if format == "markdown" {
		return buf.Bytes(), nil
	}

	return md2man.Render(buf.Bytes()), nil
}
