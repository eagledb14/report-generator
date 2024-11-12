package createform

import (
	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/templates"
)

type PortViewer struct {
	Events []*alerts.Event
}

func (p *PortViewer) CreateMarkdown() string {
	const page = `
{{range .Events}}
### [{{.Ip}}]({{.HostLink}})
{{range $key, $value := .Ports}}
{{$key}}
{{range $value}}
- [{{.Name}}](https://www.cve.org/CVERecord?id={{.Name}}) Priority: {{.Rank}}
{{end}}
{{end}}
{{end}}
`

	return templates.ExecuteText("portViewermd", page, p)
}
