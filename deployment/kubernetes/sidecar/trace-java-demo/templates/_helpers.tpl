{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "trace-java-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "trace-java-demo.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "trace-java-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "trace-java-demo.labels" -}}
helm.sh/chart: {{ include "trace-java-demo.chart" . }}
{{ include "trace-java-demo.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "trace-java-demo.labels-db" -}}
helm.sh/chart: {{ include "trace-java-demo.chart" . }}
{{ include "trace-java-demo.selectorLabels-db" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "trace-java-demo.labels-authorsapp" -}}
helm.sh/chart: {{ include "trace-java-demo.chart" . }}
{{ include "trace-java-demo.selectorLabels-authorsapp" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "trace-java-demo.labels-booksapp" -}}
helm.sh/chart: {{ include "trace-java-demo.chart" . }}
{{ include "trace-java-demo.selectorLabels-booksapp" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}


{{/*
Selector labels
*/}}
{{- define "trace-java-demo.selectorLabels" -}}
app.kubernetes.io/name: {{ include "trace-java-demo.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "trace-java-demo.selectorLabels-db" -}}
app.kubernetes.io/name: {{ include "trace-java-demo.name" . }}-db
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "trace-java-demo.selectorLabels-authorsapp" -}}
app.kubernetes.io/name: {{ include "trace-java-demo.name" . }}-authorsapp
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "trace-java-demo.selectorLabels-booksapp" -}}
app.kubernetes.io/name: {{ include "trace-java-demo.name" . }}-booksapp
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}


{{/*
Create the name of the service account to use
*/}}
{{- define "trace-java-demo.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "trace-java-demo.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
