package types

import (
	"time"
	"github.com/eagledb14/form-scanner/alerts"
)

type State struct {
    FeedEvents []*alerts.Event
    Events []*alerts.Event
    Name string
    FormType Form
    EventIndex int
}

func NewState() State {
    feedEvents := alerts.DownloadRss()

    for _, e := range feedEvents {
	time.Sleep(time.Duration(1))
	go func(e *alerts.Event) {
	    e.Load()
	}(e)
    }

    return State{
	FeedEvents: feedEvents,
    }
}

func (e *State) GetFeedEvent(index int) *alerts.Event {
    return e.FeedEvents[index + (e.EventIndex * 10)]
}
