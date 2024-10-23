package createform

import "github.com/eagledb14/form-scanner/alerts"

type OpenPort struct {
	FormNumber string
	Threat string
	Summary string
	Body string
	Tlp string
	Reference string
	Events []*alerts.Event
}
