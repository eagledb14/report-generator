package alerts

import (
	"sort"
	"strings"
)

type Credentials struct {
	Email    string
	Password bool
	Source   string
	LeakType string
	AccountType string
}

func ParseCredentialDump(passwordDump string) []Credentials {
	credentialMap := make(map[string]Credentials)

	lines := strings.Split(passwordDump, "\n")

	for i, line := range lines {
		if strings.Contains(line, "Source") {
			source := parseSource(line)
			email, pass := parseEmail(lines[i-1])
			credentialMap[email] = Credentials{
				Email:    strings.TrimLeft(email, " "),
				Password: pass,
				Source:   source,
				LeakType: "Email Address",
				AccountType: "Agency",
			}
		}
	}

	credentials := []Credentials{}
	for _, value := range credentialMap {
		credentials = append(credentials, value)
	}

	credentials = SortCreds(credentials)

	return credentials
}

func ParseOtherCreds(creds string) []Credentials {
	credList := strings.Split(creds, "\n")
	newCredentials := []Credentials{}

	for _, cred := range credList {
		credSplit := strings.Split(cred, ",")
		newCred := Credentials{}

		if len(credSplit) == 1 {
			continue
		}

		if len(credSplit) >= 5 {
			newCred.AccountType = credSplit[4]
		}
		if len(credSplit) >= 4 {
			newCred.LeakType = credSplit[3]
		}
		if len(credSplit) >= 3 {
			newCred.Source = credSplit[2]
		}
		if len(credSplit) >= 2 {
			if credSplit[1] != "" && credSplit[1] != " " {
				newCred.Password = true
			}
		}
		newCred.Email = credSplit[0]

		newCredentials = append(newCredentials, newCred)
	}

	newCredentials = SortCreds(newCredentials)

	return newCredentials
}

func SortCreds(creds []Credentials) []Credentials {
	sort.Slice(creds, func(i, j int) bool {
		return creds[i].Email < creds[j].Email
	})
	return creds
}

func parseSource(source string) string {
	parts := strings.Split(source, " ")
	_ = parts
	endIndex := len(parts)
	for i, val := range parts {
		if val == "on" {
			endIndex = i
		}
	}
	return strings.Join(parts[1:endIndex], " ")
}

func parseEmail(email string) (string, bool) {
	parts := strings.Split(email, ":")
	hasPassword := false

	if len(parts) >= 2 {
		hasPassword = true
	}
	if len(parts[0]) > 0 && parts[0][len(parts[0])-1] == '.' {
		return parts[0][:len(parts[0])-1], hasPassword
	}
	return strings.ToLower(parts[0]), hasPassword
}
