package types

import (
    "github.com/eagledb14/form-scanner/alerts"
)

type State struct {
    Feed alerts.Feed
    Events []*alerts.Event
    Name string
}

func NewState() State {
    return State{}
}
