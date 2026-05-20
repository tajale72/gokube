{{/*
Expand the name of the chart.
*/}}
{{- define "gokube.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "gokube.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gokube.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gokube.labels" -}}
helm.sh/chart: {{ include "gokube.chart" . }}
{{ include "gokube.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gokube.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gokube.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
API Service labels
*/}}
{{- define "gokube.api.labels" -}}
helm.sh/chart: {{ include "gokube.chart" . }}
app.kubernetes.io/name: {{ .Values.api.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
API Service selector labels
*/}}
{{- define "gokube.api.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.api.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
{{- end }}

{{/*
UI Service labels
*/}}
{{- define "gokube.ui.labels" -}}
helm.sh/chart: {{ include "gokube.chart" . }}
app.kubernetes.io/name: {{ .Values.ui.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: ui
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
UI Service selector labels
*/}}
{{- define "gokube.ui.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.ui.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: ui
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "gokube.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "gokube.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create image name with registry
*/}}
{{- define "gokube.api.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s:%s" .Values.global.imageRegistry .Values.api.image.repository .Values.api.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.api.image.repository .Values.api.image.tag }}
{{- end }}
{{- end }}

{{/*
Create image name with registry for UI
*/}}
{{- define "gokube.ui.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s:%s" .Values.global.imageRegistry .Values.ui.image.repository .Values.ui.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.ui.image.repository .Values.ui.image.tag }}
{{- end }}
{{- end }}

