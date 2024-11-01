package createform

import (
	"encoding/csv"
	"strings"

	"github.com/eagledb14/form-scanner/alerts"
)


func CreateCsv(query string) string {
	net :=	alerts.DownloadMatches(query)
	data := [][]string{
		{"asn", "ip", "port", "timestamp", "domains", "data", "hostnames", "isp", "org", "os", "country", "country code", "region code", "city", "product"},
	}

	for _, match := range net.Matches {
		row := []string{
			match.Asn,
			match.Ip,
			match.Port,
			match.Timestamp,
			strings.Join(match.Domains, " "),
			strings.Join(match.Hostnames, " "),
			match.Isp,
			match.Org,
			match.Os,
			match.Location.CountryName,
			match.Location.CountryCode,
			match.Location.RegionCode,
			match.Location.City,
			match.Product,
		}
		data = append(data, row)
	}

	var csvString strings.Builder
	writer := csv.NewWriter(&csvString)

	for _, record := range data {
		 writer.Write(record)
	}

	writer.Flush()

	return csvString.String()
}
