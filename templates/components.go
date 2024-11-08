package templates

import (
	"bytes"
	"html/template"
	text "text/template"

	"github.com/eagledb14/form-scanner/types"
)

func Execute(name string, t string, data interface{}) string {
	tmpl, err := template.New(name).Parse(t)
	if err != nil {
		return err.Error()
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return err.Error()
	}

	return b.String()
}

// The reason both functions are needed is because html/template sanitizes
// the html input, which is something we want, unless we already
// sanitized the input
func ExecuteText(name string, t string, data interface{}) string {
	tmpl, err := text.New(name).Parse(t)
	if err != nil {
		return err.Error()
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return err.Error()
	}

	return b.String()
}

func Banner(state *types.State) string {
	data :=	struct {
		EventIndex int
	} {
		EventIndex: state.EventIndex,
	}

	const page =  `
        <div class="heading">
            <nav style="margin: 0px 10px">
                    <ul>
						<li>
							<details class="dropdown">
								<summary role="button" class="contrast">
									Reports
								</summary>
								<ul dir="ltl">
									<li><a href="/actor">Actor</a></li>
									<li><a href="/credleak">Cred Leak</a></li>
									<li><a href="/csv">CSV</a></li>
									<li><a href="/event/page/{{.EventIndex}}">Event</a></li>
									<li><a href="/openport">Open Port</a></li>
									<li><a href="/portview">Port Viewer</a></li>
								</ul>
							</details>
						</li>
                        <li><a  role="button" class="contrast" href="/preview">Markdown Preview</a></li>
                    </ul>
            </nav>
        </div>
    `
	return Execute("banner", page, data)
}

func header() string {
	return `
        <head>
            <title>JCTF Form Generator</title>
            <script src="https://unpkg.com/htmx.org@1.9.12" integrity="sha384-ujb1lZYygJmzgSwoxRggbCHcjc0rB2XoQrxeTUQyRjrOnlCoYta87iKBWq3EsdM2" crossorigin="anonymous"></script>
	    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.blue.min.css">
	    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.colors.min.css">
            <link rel="stylesheet" type="text/css" href="/style.css">
        </head>
        `
}

func BuildPage(body string, state *types.State) string {
	data := struct {
		Header string
		Body   string
		Banner string
	}{
		Header: header(),
		Body:   body,
		Banner: Banner(state),
	}

	const page = `
        <!DOCTYPE html>
        <html lang="en">
        {{.Header}}
        <body hx-boost="true">
	    {{.Banner}}
            <div class="center">
                {{.Body}}
            </div>
        </body>
        </html>
        `

	return ExecuteText("page", page, data)
}
