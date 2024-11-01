package createform

import (
	"github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
)

type Actor struct {
	Name         string
	Alias        string
	Date         string
	Country      string
	Motivation   string
	Target       string
	Malware      string
	Reporter     string
	Confidence   string
	Exploits     string
	Summary      string
	Capabilities string
	Detection    string
	Ttps         string
	Infra        string
}

func (a *Actor) CreateMarkdown(state *types.State) string {

	const page = `
## OVERVIEW

Primary name: {{.Name}}

Alias: {{.Alias}}

First Seen Activity: {{.Date}}

Country of Origin: {{.Country}}

Motivation: {{.Motivation}}

### ATTRIBUTION ASSESSMENT

Assessment Confidence: {{.Confidence}}

Assessment Details: This malware is unique and only seen with {{.Country}}s' nation state attacks. According to {{.Reporter}}, the {{.Name}} malware family is the name given to malware developed and controlled by an intelligence directorate supporting the nation state {{.Country}}.

Exploits:

{{.Exploits}}

### MALWARE: {{.Name}}

Capabilities: {{.Capabilities}}

## ATTACK CHAIN

### Summary

{{.Summary}}

### TTPS

{{.Ttps}}

### INFRASTRUCTURE

Use of email addresses to register infrastructure beginning in [yyyy]:

{{.Infra}}

### Targeting

{{.Target}}
`


	return templates.ExecuteText("actormd", page, a)
}

