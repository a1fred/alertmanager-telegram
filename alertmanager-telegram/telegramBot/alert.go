package telegramBot

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
)

var alertHtmlTemplate = `
{{- range .Alerts }}
{{- if eq .Status "firing" }}
<b>🔥 {{ .Status | toUpper }} 🔥</b>
{{- else }}
<b>👍 {{ .Status | toUpper }} 👍</b>
{{- end }}
<strong>Labels</strong>:
{{- range $key, $value := .Labels }}
- <b>{{ $key }}</b>: {{ $value }}
{{- else }}
- No labels
{{- end }}
<strong>Annotations</strong>:
{{- range $key, $value := .Annotations }}
- <b>{{ $key }}</b>: {{ $value }}
{{- else }}
- No labels
{{- end }}
<b>Starts</b>: {{ .StartsAt | timeFormat }} <i>{{ .StartsAt | since }} ago</i>
{{- if ne .Status "firing"}}
<b>End</b>: {{ .EndsAt | timeFormat }} <i>{{ .EndsAt | since }} ago</i>
<b>Duration</b>: {{ duration .StartsAt .EndsAt }}
{{- end }}

{{- end }}
`

var messageTemplate *template.Template

func init() {
	var err error
	messageTemplate = template.New("").Funcs(template.FuncMap{
		"toUpper": strings.ToUpper,
		"timeFormat": func(t time.Time) string {
			return t.Format("Mon, 02 Jan 2006 15:04:05 MST")
		},
		"since": func(t time.Time) string {
			return time.Since(t).Round(time.Second).String()
		},
		"duration": func(start time.Time, end time.Time) string {
			return end.Sub(start).Round(time.Second).String()
		},
	})
	messageTemplate, err = messageTemplate.Parse(alertHtmlTemplate)
	if err != nil {
		panic(err)
	}
}

func FormatAlertHtml(message webhook.Message) (string, error) {
	tpl := bytes.Buffer{}

	err := messageTemplate.Execute(&tpl, message)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
