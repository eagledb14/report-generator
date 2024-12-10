package createform

import (
	"sort"

	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/templates"
)

type Osint struct {
	Name            string
	InScope         []string
	OutScope        []string
	Events          []*alerts.Event
	Url             string
	UrlIps		[]*alerts.Event
	VulnerableUrls	int
	Creds           []alerts.Credentials
	AssetSeverity   string
	AccountSeverity string
	WebsiteSeverity string
}

func (o *Osint) CreateMarkdown() string {
	for _, event := range o.Events {
		event.FilterCves()
	}
	data := struct {
		Name            string
		InScopeIps      []string
		OutScopeIps     []string
		Events          []*alerts.Event
		CveDisplay      string
		AssetSeverity   string
		AccountSeverity string
		WebsiteSeverity string

		MaxPriority []alerts.Cve

		Urls             string
		VulnerableUrls int
		UrlIpDisplay string

		Creds        []alerts.Credentials
		NumEmails    int
		NumPasswords int

		Recommendations string
		Disclaimer      string
	}{
		Name:        o.Name,
		InScopeIps:  o.InScope,
		OutScopeIps: o.OutScope,
		Events:      alerts.FilterEvents(o.Events),
		CveDisplay:  displayCves(o.Events),

		AssetSeverity:   o.AssetSeverity,
		AccountSeverity: o.AccountSeverity,
		WebsiteSeverity: o.WebsiteSeverity,

		MaxPriority: filterMaxPriority(o.Events),

		Urls:             o.Url,
		VulnerableUrls: o.VulnerableUrls,
		UrlIpDisplay: displayUrlCves(o.UrlIps, o.Url),

		Creds:        o.Creds,
		NumEmails:    len(o.Creds),
		NumPasswords: countPassowords(o.Creds),

		Recommendations: recommendations(o.Name, true),
		Disclaimer:      disclaimer(o.Name),
	}

	const page = `
## Overall Risk Exposure Ratings

The information below shows the numbers of issues identified in different categories. Exposures are classified according to severity as Critical, High, Moderate, or Low. This reflects the likely impact of each issue for a typical organization. All information provided in this report was gathered entirely passively, no interaction with {{.Name}} or assets owned by the {{.Name}} to include scanning, crawling or active enumeration was done to gather the information in this report. 

| Exposure | Description | Severity | Count |
|---|---|---|---|
| Detect And Identify Vulnerable External Assets | Identifying Vulnerable External Assets Using Open-source Tools | {{.AssetSeverity}} | {{len .Events}} Vulnerable External Assets |
| Identify User Accounts Through Open-sources | Using Open-source Tools to Identify Exposed Accounts | {{.AccountSeverity}} | {{.NumEmails}} Emails, {{.NumPasswords}} Passwords |
| Discover Vulnerable Websites | Using Open-source Tools to Discover Domains, Subdomains and Hostnames | {{.WebsiteSeverity}} | {{.VulnerableUrls}} Vulnerable websites identified |


## External Asset Discovery

The initial phase of the passive review is External Asset Discovery. In this section, the North Carolina National Guard (NCNG) searches through open-source databases, search-engines, and Pastebin to discover external assets owned by {{.Name}}. All of the information within this section is intended to help your organization gain better visibility and exposed asset awareness of external assets.

### External IPs Provided Within Scope

{{range .InScopeIps}}
- {{.}}
{{end}}

### External IPs Indentified Outside Of Scope

{{range .OutScopeIps}}
- {{.}}
{{end}}

The IPs provided above are associated with {{.Name}}’s domain on open-source DNS record websites. Provided solely for {{.Name}}’s awareness.

---

## Identifying Vulnerable External Devices with Shodan

Shodan.io is an open-source search engine that is designed to gather information about internet-connected devices and systems. The NCNG searched Shodan’s public database for any assets owned by {{.Name}} using the CIDR Blocks or IP addresses provided within scope and identified through asset discovery. 

The table below uses the Exploit Prediction Scoring System (EPSS) and Common Vulnerability Scoring System (CVSS) to measure vulnerabilities. EPSS produces prediction scores between 0 and 1 (0 and 100%) where higher scores suggest probability of exploit and CVSS rates the severity of a vulnerability. Vulnerabilities are prioritized in order from 1 to 4, 1+ being the most severe and 4 being the least severe.

Within the list of IP addresses above, {{len .Events}} vulnerable asset(s) are indexed by Shodan.

{{.CveDisplay}}

{{if gt (len .Events) 0}}
### Impact to Agency (External Asset Vulnerabilities)

{{if gt (len .MaxPriority) 0}}
It is essential to recognize that external assets, which {{.Name}} may not be fully aware of, could pose significant risks. These risks might encompass unpatched software, misconfigurations, exposed sensitive data, or critical vulnerabilities, such as the {{len .MaxPriority}} number of Priority 1 CVEs expanded upon below. 

{{range .MaxPriority}}
- **{{.Name}}** {{.Summary}}
{{end}}

Such vulnerabilities emphasize the importance of proactive asset discovery, patch management, and security measures to safeguard {{.Name}} from these vulnerabilities. 
{{else}}
It is essential to recognize that external assets, which the {{.Name}} may not be fully aware of, could pose significant risks. These risks might encompass unpatched software, misconfigurations, exposed sensitive data, and unidentified vulnerabilities. Such risks emphasize the importance of proactive asset discovery, patch management, and security measures to safeguard the {{.Name}} from these vulnerabilities.

Such vulnerabilities emphasize the importance of proactive asset discovery, patch management, and security measures to safeguard {{.Name}} from these vulnerabilities.{{end}}{{end}}

## Vulnerable Websites
The NCNG Searched open-source databases and dark web bug bounty markets for vulnerabilities associated with {{.Urls}} and found {{if eq .VulnerableUrls 1}}1 finding{{else if gt .VulnerableUrls 0}}{{.VulnerableUrls}} findings{{else}}no issues{{end}}. 

### Impact to Agency (Vulnerable Websites)
{{.UrlIpDisplay}}


{{if gt (len .Creds) 0}}
---
## Exposed Credentials

The NCNG employs open-source tools to proactively search for instances of credential exposure. These tools are specifically designed to identify leaked credentials by crawling through data breaches, popular PasteBin websites, such as ihavebeenpwned and other common public exposure locations where compromised accounts may be found. By leveraging these open-source tools, the NCNG aims to detect any instances where sensitive information, such as usernames and passwords, have been compromised and/or publicly disclosed. This proactive approach helps mitigate the risk of unauthorized access to {{.Name}}’s systems or personnel accounts.

The tools systematically scan and analyze data breaches, which are incidents where large amounts of confidential information are illegally obtained and made publicly available. Additionally, they monitor platforms like PasteBins, which are websites commonly used for sharing text-based content, including compromised account credentials. Through the use of these open-source tools, the NCNG can swiftly identify compromised accounts associated with {{.Name}} personnel or systems. This allows for immediate action, such as resetting passwords, notifying affected individuals, and implementing additional security measures to prevent unauthorized access or potential misuse of compromised credentials.

By actively searching for credential exposure using open-source tools, {{.Name}} demonstrates a proactive and vigilant approach to safeguarding their digital assets and protecting sensitive information from falling into the wrong hands.
	
| Email | Leak Type | Source | Compromised Account Type |
|---|---|---|---|{{range .Creds}}
| {{.Email}} | {{.LeakType}} | {{.Source}} | {{.AccountType}} |{{end}}

### Impact to Agency (Exposed Credentials)
It is important to understand the risks that exposed credentials pose to {{.Name}}’s security posture, as they can grant unauthorized access to sensitive systems, data, and infrastructure. When credentials are compromised, malicious actors can utilize this information for a variety of harmful activities, such as obtaining initial access, laterally moving within the network, exiling sensitive data, or launching further attacks. By actively uncovering these exposed credentials, the NCNG aims to raise awareness within {{.Name}} and facilitate the implementation of appropriate security measures to mitigate potential threats.

While the discovery of leaked credentials in this report raises valid concerns, it is important to emphasize that their presence does not necessarily imply that a user's domain account has been compromised. Credentials can be exposed through various channels, including breaches of unrelated external services, and may not directly impact {{.Name}}’s internal systems. Additionally, if robust security measures like multi-factor authentication (MFA) are in place, the risk of unauthorized access remains significantly reduced, even if the credentials have been exposed. By identifying these leaks, The NCNG aims to proactively address potential risks and ensure that appropriate security measures are implemented.

{{end}}

---

{{.Recommendations}}

By implementing these recommendations, {{.Name}} can enhance its security posture and mitigate potential risks associated with external assets. The NCNG is committed to assisting and supporting {{.Name}}  in strengthening their cybersecurity defenses and ensuring the protection of their critical assets.

---

{{.Disclaimer}}
`

	return templates.ExecuteText("osintmd", page, data)
}

func displayCves(events []*alerts.Event) string {

	data := struct {
		Events []*alerts.Event
	}{
		Events: events,
	}

	const page = `
{{if eq (len .Events) 0}}{{else if eq (len .Events) 1}}{{$event := index .Events 0}}
The external IP; “{{$event.Ip}}” is tagged with vulnerabilities on Shodan. A table of the vulnerabilities for this IP is found below.
{{else}}
The external IPs; {{range .Events}}“{{.Ip}}”, {{end}}are tagged with vulnerabilities on Shodan. Tables of the vulnerabilities for these IPs are found below.{{end}}

{{if eq (len .Events) 0}}{{else}}{{range .Events}}{{if gt (len .Ports) 0}}

**[{{.Ip}}]**

| CVE-ID | PRIORITY | EPSS | CVSS | VERSION | SEVERITY | CISA_KEV | VENDOR | PRODUCT |
|---|---|---|---|---|---|---|---|---|{{range $key, $cve := .Ports}}{{range $cve}}
| {{.Name}} | Priority {{.Rank}} | {{.Epss}} | {{.Cvss}} | {{.Version}} | {{.Severity}} | {{.Kev}} | {{.Vendor}} | {{.Product}} |{{end}}{{end}}{{end}}{{end}}{{end}}`

	return templates.Execute("displayCves", page, data)
}

func displayUrlCves(events []*alerts.Event, url string) string {
	uniqueCve := make(map[string]alerts.Cve)
	for _, e := range events {
		for _, cve := range e.Ports {
			for _, c := range cve {
				uniqueCve[c.Name] = c
			}
		}
	}

	cves := []alerts.Cve{}
	for _, cve := range uniqueCve {
		cves = append(cves, cve)
	}

	sort.Slice(cves, func(i, j int) bool {
		return cves[i].Rank < cves[j].Rank
	})

	data := struct {
		Events []*alerts.Event
		Cves []alerts.Cve
		Url string
	} {
		Events: events,
		Cves: cves,
		Url: url,
	}
	const page = `
{{range .Events}}
**{{.Ip}}: {{$.Url}}**
{{end}}
<br>

{{range .Cves}}
- {{.Name}}: Priority {{.Rank}}
	- {{.Summary}}
{{end}}
`

	return templates.Execute("displayUrlCves", page, data)
}

func recommendations(name string, creds bool) string {
	data := struct {
		Name  string
		Creds bool
	}{
		Name:  name,
		Creds: creds,
	}

	const page = `
## Recommendations:

Based on the findings from the External Asset Discovery phase, the NCNG strongly recommends the following actions to {{.Name}}:

1. Regular Vulnerability Assessments: Conduct comprehensive and periodic vulnerability assessments of all identified external assets. This will help identify and prioritize vulnerabilities that need to be addressed promptly.
2. Configuration Reviews: Perform regular reviews of configurations for all external assets to ensure they align with industry best practices and security standards. Misconfigurations can often lead to security weaknesses that attackers can exploit, so maintaining secure configurations is crucial.
3. Access Control and Authentication: Strengthen access controls for external assets by implementing strong authentication mechanisms, such as multifactor authentication (MFA), and regularly reviewing and revoking unnecessary privileges. This will help protect against unauthorized access attempts.
4. Monitoring and Incident Response: Establish robust monitoring capabilities to detect and respond to any suspicious activity or potential breaches related to external assets. Implement an incident response plan that outlines the necessary steps to be taken in the event of a security incident.
5. Employee Awareness and Training: Conduct regular cybersecurity awareness training for employees to educate them about the risks associated with external assets and how to follow best practices for security. This will help foster a culture of cybersecurity within {{.Name}} and empower employees to contribute to the overall protection of external assets.
{{if .Creds}}6. Account Usage: Employees are strongly advised to refrain from using their work email addresses for personal service sign-ups to reduce the risk of credentials being exposed through third-party breaches.{{end}}
`
	return templates.ExecuteText("recommendations", page, data)
}

func disclaimer(name string) string {
	return `CUI / / FOR OFFICIAL USE ONLY. This document is the exclusive property of ` + name + `. It contains proprietary and confidential information and may not be duplicated, redistributed, or used, in whole or in part, in any form, without consent of ` + name + ` and North Carolina Cyber Security Response Force. Restricted / Confidential per N.C.G.S. § 132-6.1(c)`
}

func filterMaxPriority(events []*alerts.Event) []alerts.Cve {
	maxCves := []alerts.Cve{}

	for _, event := range events {
		for _, cves := range event.Ports {
			for _, cve := range cves {
				if cve.Rank < 2 {
					maxCves = append(maxCves, cve)
				}
			}
		}
	}

	return maxCves
}

func countPassowords(creds []alerts.Credentials) int {
	numPasswords := 0

	for _, cred := range creds {
		if cred.Password {
			numPasswords += 1
		}
	}

	return numPasswords
}
