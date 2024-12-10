package alerts

import (
	"fmt"
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

	emails := []string{}
	for _, line := range lines {
		if strings.Contains(line, "Credential leak") {
			emails = []string{}

			emailSplit := strings.Split(line, " ")

			for _, text := range emailSplit {
				if strings.Contains(text, "@") {
					text = strings.ReplaceAll(text, ",", "")
					emails = append(emails, text)
				}
			}

			for _, email := range emails {
				credentialMap[email] = Credentials{
					Email: email,
					LeakType: "Email Address",
					AccountType: "Agency",
				}
			}
		}

		if strings.Contains(line, "Source") {
			source := parseSource(line)
			for _, email := range emails {
				cred, _ := credentialMap[email]
				cred.Source = source
				credentialMap[email] = cred
			}
		}

		if strings.HasPrefix(line, "    ") {
			for _, email := range emails {
				cred, _ := credentialMap[email]
				cred.Password = hasPassword(line)
				credentialMap[email] = cred
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

func hasPassword(line string) bool {
	parts := strings.Split(line, ":")

	return len(parts) >= 2
}

