package alerts

type Cve struct {
	Name string
	Rank int
}

func NewCve(name string, vuln Vuln) Cve {
	var cvss float32 = 0.0
	if vuln.Cvss > vuln.CvssV2 {
		cvss = vuln.Cvss
	} else {
		cvss = vuln.CvssV2
	}
	rank := rankCve(cvss, vuln.Epss, vuln.Kev)

	return Cve{
		Name: name,
		Rank: rank,
	}
}

func rankCve(cvss float32, epss float32, kev bool) int {
	var cvssScore float32 = 6.0
	var epssScore float32 = 0.2

	if kev {
		return 0
	} else if cvss >= cvssScore {
		if epss >= epssScore {
			return 1
		} else {
			return 2
		}
	} else {
		if epss >= epssScore {
			return 3
		} else {
			return 4
		}
	}
}
