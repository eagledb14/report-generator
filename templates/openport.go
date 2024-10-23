package templates

import (
	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/types"
)

// or I might pass in state to this
func OpenPortDownload() string {
	data := struct {
	}{}

	const page = `
        <h1>Open Port</h1>
		<article>
			<form hx-post="/openportform" hx-target="body">
				<fieldset>
						<label>
							Organization Name
							<input name="orgName"/>
						</label>
						<label>
							IP Addresses
							<input name="ipAddress" />
						</label>

						<div class="grid">
							<input type="submit">
							<input type="reset">
						</div>
				</fieldset>
			</form>
		</article>
        `

	return Execute("openport", page, data)
}

func OpenPortForm(name string, e []*alerts.Event) string {
	data := struct {
		Name string
		Events []*alerts.Event
		Form string
	}{
		Name: name,
		Events: e,
		Form: form(types.Open, e),
	}

	const page = `
		<button hx-put="/openport" hx-target="body"><</button>
        <h1>{{.Name}}</h1>
		{{range .Events}}
			<article>
				<header><h3>{{.Ip}}<h3></header>
				{{range $key, $value := .Ports}}
					<h4>{{$key}}</h4>
					{{range $value}}
						<small>{{.Name}}: Priority {{.Rank}}</small>
						<br>
					{{end}}
					<hr>
				{{end}}
			</article>
		{{end}}
		{{.Form}}

        `

	return ExecuteText("openport", page, data)
}
