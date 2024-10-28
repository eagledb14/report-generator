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
			<form hx-post="/openport/form" hx-target="body" hx-indicator="#load">
				<fieldset>
						<label>
							Organization Name
							<input name="orgName"/>
						</label>
						<label>
							IP Addresses
							<input name="ipAddress" />
						</label>

						<div id="load" class="htmx-indicator center" aria-busy="true">Loading...</div>
						<div class="grid">
							<input type="submit" value="Submit"/>
							<input type="reset"/>
						</div>
				</fieldset>
			</form>
		</article>
        `

	return Execute("openport", page, data)
}

func OpenPortForm(form types.Form, name string, e []*alerts.Event) string {
	data := struct {
		Name string
		Events []*alerts.Event
		Form string
		FormName string
	}{
		Name: name,
		Events: e,
		Form: getForm(form, name, e, "/openport"),
		FormName: types.FormName[form],
	}

	const page = `
		<button hx-put="/openport" hx-target="body"><</button>
        <h1>{{.Name}}</h1>
		{{range .Events}}
			<article>
				<header>
					<h3>{{.Ip}}</h3> 
					<br>
					<small><a href="{{.HostLink}}" target=_blank>Host Link</a></small>
				</header>
				{{if eq (len .Ports) 0}}
					<h4> No Available Information</h4>
				{{end}}
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
		<hr>
		<div class="grid">
			<button hx-get="/openport/port" hx-target="body">Open Port</button>
			<button hx-get="/openport/eol" hx-target="body">End of Life</button>
			<button hx-get="/openport/login" hx-target="body">Login Pages</button>
		</div>
		<h3>{{$.FormName}}</h3>
		{{.Form}}
        `

	return ExecuteText("openport", page, data)
}
