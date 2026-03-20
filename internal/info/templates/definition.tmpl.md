# Template: {{ .TemplateName }}
# Definition: {{ .DefinitionName }} ({{ .Source }})

## Variables
| name | type | required | value |
|------|------|----------|-------|
{{- if .Variables }}
{{- range .Variables }}
| {{ .Name }} | {{ .Type }} | {{ boolCell .Required }} | {{ valueCell .Value .Default }} |
{{- end }}
{{- else }}
| (none) |  |  |  |
{{- end }}

## Files
{{- if .Targets }}
{{- range .Targets }}
- {{ . }}
{{- end }}
{{- else }}
- (none)
{{- end }}
