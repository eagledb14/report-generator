package createform

import (
	"github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
)


type CredLeak struct {
	OrgName string
	FormNumber string
	VictimOrg string
	Password string
	UserPass string
	AddInfo string
	Reference string
	Tlp bool
}

func (c *CredLeak) CreateMarkdown(state *types.State) string {
	state.AlertId = getAlertId(c.FormNumber)

	data := struct {
		Name string
		AlertId string
		Password string
		UserPass string
		AddInfo string
		Reference string
		Footer string
	} {
		Name: c.OrgName,
		AlertId: getAlertId(c.FormNumber),
		Password: c.Password,
		UserPass: c.UserPass,
		AddInfo: c.AddInfo,
		Reference: c.Reference,
		Footer: footer(c.Tlp),
	}

	const page = `
## {{.Name}}

ALERT ID: {{.AlertId}}

THREAT TYPE: Credential Exposure (Mitre ATT&CK ID: T1078)

### SUMMARY

The North Carolina National Guard Cyber Security Response Force identified leaked credentials within the {{.Name}} domain. We discovered these credentials due to routine threat monitoring on the Dark Web. We do not know if these credentials can access {{.Name}} resources or are valid credentials, but we want to share them for your situational awareness and further investigation.


{{.Password}}

{{.UserPass}}

{{.AddInfo}}

### REFERENCES
{{.Reference}}

---

## Recommended Next Steps
Credential theft is one of the most common avenues of initial compromise in cyber attacks such as ransomware. The CSRF recommends that you validate that these accounts are current and if so, perform a password reset of all compromised accounts. If any of the credentials have been compromised via info-stealer malware, the malware must be removed first. 

If these credentials are not current:
1. Set alerts for failed logins from this account in your security tools or Security information and event management (SIEM)
2. Inform the victim and offer advice on how to secure personal devices such as installing anti-virus, using Multi-Factor Authentication, having good password management, and updating both the computer and software programs on it. Resource: [https://www.cisa.gov/uscert/ncas/tips/ST04-003](https://www.cisa.gov/uscert/ncas/tips/ST04-003)

If these credentials are current:
1. Implement an immediate password reset for that account and consider a domain-wide reset
2. Utilize Multi-Factor Authentication if not done already 
3. Set alerts for this account in your security tools or SIEM
4. Investigate your enterprise to determine if activity from this account is a security incident
5. Inform the victim and offer advice on how to secure personal devices such as installing anti-virus, using Multi-Factor Authentication, having good password management, and both updating the computer and software programs on it 
6. Inform executive leadership and your security officer about the compromise and include credential theft awareness training in your next information security training.

If you identify any of these credentials as a security incident, and you require assistance, please contact the North Carolina National Guard CSRF for cyber incident response support.

If any of these credentials are a false positive please consider letting us know so that we know it is not associated with a cyber threat actor.

---

## Additional Resources
- [Good Security Habits](https://us-cert.cisa.gov/ncas/tips/ST04-003)
- [Understanding Anti-Virus Software](https://us-cert.cisa.gov/ncas/tips/ST04-005)
- [Good Security Habits](https://us-cert.cisa.gov/ncas/tips/ST04-003)
- [Understanding Anti-Virus Software](https://us-cert.cisa.gov/ncas/tips/ST04-005)
- [Understanding Patches and Software Updates](https://us-cert.cisa.gov/ncas/tips/ST04-006)
- [Choosing and Protecting Passwords](https://us-cert.cisa.gov/ncas/tips/ST04-002)
- [SMB Security Best Practices](https://us-cert.cisa.gov/ncas/current-activity/2017/01/16/SMB-Security-Best-Practices)
- [Rising Ransomware Threat to Operational Technology Assets](https://www.cisa.gov/sites/default/files/publications/CISA_Fact_Sheet-Rising_Ransomware_Threat_to_OT_Assets_508C.pdf)

---

{{.Footer}}
`

	return templates.ExecuteText("credleakmd", page, data)
}

