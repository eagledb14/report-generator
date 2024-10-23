package types

import (
    "github.com/eagledb14/form-scanner/alerts"
)

type State struct {
    feed alerts.Feed
    event []*alerts.Event
}
