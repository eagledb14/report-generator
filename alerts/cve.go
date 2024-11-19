package alerts

import (
	"strings"
)

type Cve struct {
	Name string
	Summary string
	Rank int
	Epss float32
	Cvss float32
	Kev bool
	Version string
	Severity string
	Vendor string
	Product string
}

func NewCve(name string, vuln Vuln, cpe []string) Cve {
	newCve := Cve{}
	var cvss float32 = 0.0
	
	if vuln.Cvss > vuln.CvssV2 {
		cvss = vuln.Cvss
		newCve.Version = "CVSS 1.0"
	} else {
		cvss = vuln.CvssV2
		newCve.Version = "CVSS 2.0"
	}
	rank, severity := rankCve(cvss, vuln.Epss, vuln.Kev)
	vendor, product := getVendorProduct(cpe)

	newCve.Name = strings.TrimLeft(name, " ")
	newCve.Summary = vuln.Summary
	newCve.Rank = rank
	newCve.Epss = vuln.Epss
	newCve.Cvss = cvss
	newCve.Kev = vuln.Kev
	newCve.Severity = severity
	newCve.Vendor = vendor
	newCve.Product = product

	return newCve
}

func rankCve(cvss float32, epss float32, kev bool) (int, string) {
	var cvssScore float32 = 6.0
	var epssScore float32 = 0.2

	if kev {
		return 0, "HIGH"
	} else if cvss >= cvssScore {
		if epss >= epssScore {
			return 1, "HIGH"
		} else {
			return 2, "MODERATE"
		}
	} else {
		if epss >= epssScore {
			return 3, "MODERATE"
		} else {
			return 4, "LOW"
		}
	}
}

func getVendorProduct(cpes []string) (string, string) {
	for _, cpe := range cpes {
		cpe_split := strings.Split(cpe, ":")
		if len(cpe_split) > 4 {
			return cpe_split[3], cpe_split[4]
		}
	}

	return "", ""
}
