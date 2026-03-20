# Templates
{{- if .InProject }}
- {{ .CurrentTemplate }} (template)

# Definitions
{{- if .LocalDefinitions }}
{{- range .LocalDefinitions }}
- {{ . }} (local)
{{- end }}
{{- else }}
- (none)
{{- end }}
{{- else }}
{{- range .TemplateNames }}
- {{ . }}
{{- end }}
{{- end }}
