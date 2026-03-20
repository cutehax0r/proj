# Template: {{ .TemplateName }}

## New
{{- if .NewDefs }}
{{- range .NewDefs }}
- {{ .Name }} ({{ .Source }})
{{- end }}
{{- else }}
- (none)
{{- end }}

## Add
{{- if .AddDefs }}
{{- range .AddDefs }}
- {{ .Name }} ({{ .Source }})
{{- end }}
{{- else }}
- (none)
{{- end }}
