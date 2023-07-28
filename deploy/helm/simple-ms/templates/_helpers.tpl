{{- define "app.name" -}}
{{- .Chart.Name }}
{{- end }}

{{- define "app.version" -}}
{{- default .Chart.AppVersion .Values.image.simpleMs.tag }}
{{- end }}

{{- define "app.schema.version" -}}
{{- .Values.image.pgSchema.tag }}
{{- end }}

{{- define "app.labels" -}}
app.kubernetes.io/name: {{ include "app.name" . | quote }}
app.kubernetes.io/version: {{ include "app.version" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}

{{- define "pgSchema.labels" -}}
{{ include "app.labels" . }}
app.kubernetes.io/schema-version: {{ include "app.schema.version" . | quote }}
{{- end }}

{{- define "pg.migrator.config.name" -}}
{{ .Release.Name }}-postgres-migrator-config
{{- end }}

{{- define "pg.migrator.secret.name" -}}
{{ .Release.Name }}-postgres-migrator-secret
{{- end }}

{{- define "pg.database.service.name" -}}
{{ .Chart.Name }}-postgres-db-service
{{- end }}

{{- define "app.service.name" -}}
{{ .Release.Name }}-serive
{{- end }}

{{- define "app.pod.exposedPort" -}}
8000
{{- end }}

{{- define "app.config.name" -}}
{{ .Release.Name }}-app-config
{{- end }}

{{- define "app.secret.name" -}}
{{ .Release.Name }}-app-secret
{{- end }}