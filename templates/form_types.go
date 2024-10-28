package templates

import (
	"strings"

	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/types"
)


func getForm(formType types.Form, name string, events []*alerts.Event, endpoint string) string {
	summary := ""
	body := ""

	// make a match on which type is passed int
	switch formType {
	case types.Open:
		summary = openPortSummary(name, events)
		body = openPortBody(name, events)
	case types.EOL:
		summary = endOfLifeSummary(name, events)
		body = openPortBody(name, events)
	case types.Login:
		summary = loginPageSummary(name)
		body = loginPageBody(name, events)
	default:
		summary = types.FormName[formType]
		body = types.FormName[formType]
	}

	data := struct {
		Summary string
		Body string
		Endpoint string
	} {
		Summary: summary,
		Body: body,
		Endpoint: endpoint,
	}

	const page = `
	<article>
		<form hx-post="{{.Endpoint}}" hx-target="body">
			<fieldset>
                    <label>
                        Form Number
                        <input name="formNumber"/>
                    </label>

					<label>
						Threat Type
						<input name="threat" value="T1133 External Remote Services"/>
					</label>

					<label>
						Summary Paragraph
						<textarea name="summary">{{.Summary}}</textarea>
					</label>

					<label>
						Body Paragraph
						<textarea name="body">{{.Body}}</textarea>
					</label>

					<label>
						Additional References
						<textarea name="reference"></textarea>
					</label>

					<label>TLP Alert</label>
					<label>
						<input type="radio" value="amber" name="tlp" checked/>
						Amber
					</label>
					<label>
						<input type="radio" value="green" name="tlp"/>
						Green
					</label>

					<hr>
					<div class="grid">
						<input type="submit" value="Submit" onclick="window.scrollTo(0, 0);">
						<input type="reset">
					</div>
			</fieldset>
		</form>
	</article>
	`

	return Execute("form", page, data)
}

func openPortSummary(name string, events []*alerts.Event) string {
	cves := false
	outer: for _, e := range events {
		for _, cve := range e.Ports {
			if len(cve) > 0 {
				cves = true
				break outer
			}
		}
	}

    if cves {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the " + name + " domain is publicly exposed to the internet via several risky open ports and CVEs of concern."
    } else {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the " + name + " domain is publicly exposed to the internet via several risky open ports."
    }
}

func openPortBody(name string, events []*alerts.Event) string {
	ips := strings.Builder{}

	for _, e := range events {
		ips.Write([]byte(e.Ip + ", "))
	}
	ipString := ""

	if len(ips.String()) > 0 {

		ipString = ips.String()[:len(ips.String()) - 2]
	}

	return "A threat actor may have an easier pathway to conducting a cyber attack or cyber espionage against your organization based on your current configuration. We encourage " + name + " to review the infrastructure at the following IP addresses: " + ipString + " and evaluate the risk of leaving them in their current state. We also encourage " + name + " to search for indicators of unauthorized access because threat actors exploit this configuration often for initial access."
}

func endOfLifeSummary(name string, events []*alerts.Event) string {
	cves := false
	outer: for _, e := range events {
		for _, cve := range e.Ports {
			if len(cve) > 0 {
				cves = true
				break outer
			}
		}
	}

    if cves {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the " + name + " domain is publicly exposing end of life infrastructure via several risky open ports and CVEs of concern."
    } else {
        return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the " + name + " domain is publicly exposing end of life infrastructure via several risky open ports of concern."
    }
}

func loginPageSummary(name string) string {
    return "The North Carolina National Guard Cyber Security Response Force (NCNG CSRF) received an alert indicating the " + name + " domain is publicly exposing risky login pages to the internet."
}

func loginPageBody(name string, events []*alerts.Event) string {
	ips := strings.Builder{}
	for _, e := range events {
		ips.Write([]byte(e.Ip + ", "))
	}
	ipString := ""

	if len(ips.String()) > 0 {

		ipString = ips.String()[:len(ips.String()) - 2]
	}

    return "A threat actor may have an easier pathway to conducting a cyber attack or cyber espionage against your organization based on your current configuration, through repeated login attempts against possible weak user login credentials. We encourage " + name + " to review the infrastructure at the following IP addresses: " + ipString + " and evaluate the risk of leaving them in their current state. We also encourage " + name + " to search for indicators of unauthorized access because threat actors exploit this configuration often for initial access."
}
