package createform

import (
	"strings"

	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
	// "strings"
	// "github.com/eagledb14/form-scanner/types"
)

type OpenPort struct {
	OrgName    string
	FormNumber string
	Threat     string
	Summary    string
	Body       string
	Tlp        bool
	Reference  string
	Events     []*alerts.Event
}

func (o *OpenPort) CreateMarkdown(state *types.State) string {
	state.AlertId = getAlertId(o.FormNumber)
	data := struct {
		Name       string
		AlertId    string
		ThreatType string
		Summary    string
		Body       string
		Events     string
		PriorityKey string
		Mitigations string
		Footer string
	}{
		Name:       o.OrgName,
		AlertId:    getAlertId(o.FormNumber),
		ThreatType: o.Threat,
		Summary:    o.Summary,
		Body:       o.Body,
		Events:     getEventsString(o.Events),
		PriorityKey: cvePriorityKey(),
		Mitigations: mitigations(),
		Footer: footer(o.Tlp),
	}

	const page = `
## {{.Name}} 

ALERT ID: {{.AlertId}}

THREAT TYPE: {{.ThreatType}}
	
### SUMMARY
{{.Summary}}

{{.Body}}

---

## IP and Port Assessment 

{{.Events}}

---

{{.PriorityKey}}

---

{{.Mitigations}}

---

{{.Footer}}
`
	return templates.ExecuteText("openportmd", page, data)
}

func getEventsString(events []*alerts.Event) string {
	data := struct {
		Events []*alerts.Event
	}{
		Events: events,
	}

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

	return templates.ExecuteText("eventsString", page, data)
}

func getSourcesString(events []*alerts.Event) string {
	builder := strings.Builder{}

	for _, event := range events {
		builder.WriteString(event.HostLink + "<br>")
	}

	return builder.String()
}

func cvePriorityKey() string {
	return `
## CVE Priority Key
[EPSS User Guide (first.org)](https://www.first.org/epss/user-guide)

[CVSS, EPSS, and CISA's Known Exploited Vulnerabilities (Github.com)](https://github.com/eagledb14/CVE_Prioritizer/tree/main?tab=readme-ov-file#our-approach)

<small>Priority severity is ranked 0 (highest) to 4 (lowest)</small>

### Priority 0
CISA (Cybersecurity and Infrastructure Security Agency) has declared this vulnerability as being a known exploited vulnerability. Should be taken as highest priority and addressed immediately.

### Priority 1
The most critical kinds of vulnerabilities which are more likely to be exploited, and could fully compromise the information system. They should be patched first.

### Priority 2
May severely impact the system, are much less likely to be exploited, relative to others, but should still be watched in the event that the threat landscape changes.

### Priority 3
May be more likely to be exploited, but, on their own, would not critically impact the information system.

### Priority 4
Vulnerabilities may be more likely to be exploited on their own, but would not critically impact the system.
`
}

func mitigations() string {
	return `
## Mitigations
The NCNG CSRF recommends that network defenders apply the following industry best practices to reduce the risk of compromise by ransomware attacks.

#### Preparing for Cyber Incidents 
- **Maintain offline backups of data**, and regularly maintain backup and restoration. By instituting this practice, the organization ensures they will not be severely interrupted, and that backup data will be accessible when it is needed.

- **Ensure all backup data is encrypted**, immutable (that is, cannot be altered or deleted), and covers the entire organization’s data infrastructure. Ensure your backup data is not already infected. 
- **Review the security posture of third-party vendors and those interconnected with your organization.** Ensure all connections between third-party vendors and outside software or hardware are monitored and reviewed for suspicious activity. 
- **Implement listing policies for applications and remote access that only allow systems to execute known and permitted programs** under an established security policy. 
- **Document and monitor external remote connections.** Organizations should document approved solutions for remote management and maintenance, and immediately investigate if an unapproved solution is installed on a workstation. 
- **Implement a recovery plan** to maintain and retain multiple copies of sensitive or proprietary data and servers in a physically separate, segmented, and secure location (that is, a hard drive, other storage device, or the cloud).

#### Identity and Access Management
- **Require all accounts with password logins (for example, service account, admin accounts, and domain admin accounts) to comply with National Institute of Standards and Technology (NIST) standards for developing and managing password policies.**
	- Use longer passwords consisting of at least 8 characters and no more than 64 characters in length
	- Use MFA for all accounts
	- Store passwords in hashed format using industry recognized password managers
	- Add password user “salts” to shared login credentials
	- Avoid reusing passwords;
	- Implement multiple failed login attempt account lockouts
	- Disable password “hints”
	- Refrain from requiring password changes more frequently than once per year unless a password is known or suspected to be compromised. Note: NIST guidance suggests favoring longer passwords instead of requiring regular and frequent password resets. Frequent password resets are more likely to result in users developing password “patterns” cyber criminals can easily decipher.
	- Require administrator credentials to install software. 

- **Require phishing-resistant multifactor authentication** for all services to the extent possible, particularly for webmail, virtual private networks, and accounts that access critical systems. 
- **Review domain controllers, servers, workstations, and active directories** for new and/or unrecognized accounts. 
- **Audit user accounts** with administrative privileges and configure access controls according to the principle of least privilege. 
- **Implement time-based access for accounts set at the admin level and higher.** For example, the Just-in-Time (JIT) access method provisions privileged access when needed and can support enforcement of the principle of least privilege (as well as the Zero Trust model). This is a process where a network-wide policy is set in place to automatically disable admin accounts at the Active Directory level when the account is not in direct need. Individual users may submit their requests through an automated process that grants them access to a specified system for a set timeframe when they need to support the completion of a certain task. 


#### Protective Controls and Architecture 
- **Segment networks** to prevent the spread of ransomware. Network segmentation can help prevent the spread of ransomware by controlling traffic flows between—and access to—various subnetworks and by restricting adversary lateral movement. 

- **Identify, detect, and investigate abnormal activity and potential traversal of the indicated ransomware with a networking monitoring tool.** To aid in detecting the ransomware, implement a tool that logs and reports all network traffic, including lateral movement activity on a network. Endpoint detection and response (EDR) tools are particularly useful for detecting lateral connections as they have insight into common and uncommon network connections for each host. 
- **Install, regularly update, and enable real time detection for antivirus software** on all hosts. 
- **Secure and closely monitor remote desktop protocol (RDP) use.** Limit access to resources over internal networks, especially by restricting RDP and using virtual desktop infrastructure. If RDP is deemed operationally necessary, restrict the originating sources and require MFA to mitigate credential theft and reuse. If RDP must be available externally, use a VPN, virtual desktop infrastructure, or other means to authenticate and secure the connection before allowing RDP to connect to internal devices. Monitor remote access/RDP logs, enforce account lockouts after a specified number of attempts to block brute force campaigns, log RDP login attempts, and disable unused remote access/RDP ports. 

#### Vulnerability and Configuration Management
- **Keep all operating systems, software, and firmware up to date.** Timely patching is one of the most efficient and cost-effective steps an organization can take to minimize its exposure to cybersecurity threats. Organizations should prioritize patching of vulnerabilities on CISA’s Known Exploited Vulnerabilities catalog. 

- **Do not publicly expose risky ports such as RDP Port 3389, Secure Shell (SSH) Port 22, and Server Message Block (SMB) Port 445.** 
- **Consider adding an email banner to emails** received from outside your organization. 
- **Configure enterprise email solution to disable hyperlinks and pictures** in received emails. 
- **Disable command-line and scripting activities and permissions.** Privilege escalation and lateral movement often depend on software utilities running from the command line. If threat actors are not able to run these tools, they will have difficulty escalating privileges and/or moving laterally. 
- **Ensure devices are properly configured and that security features are enabled.**

~<small>The NCNG, CISA, FBI, and NSA strongly discourage paying a ransom to criminal actors. Paying a ransom may embolden adversaries to target additional organizations, encourage other criminal actors to engage in the distribution of ransomware, and/or may fund illicit activities. Paying the ransom also does not guarantee that a victim’s files will be recovered. Additionally, your organization maybe barred from interacting and/or paying a ransom per N.C.G.S § 143-800 and § 143B-1320.</small>~
`
}
