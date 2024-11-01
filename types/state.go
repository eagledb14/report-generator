package types

import (
	"time"

	"github.com/eagledb14/form-scanner/alerts"
)

type State struct {
    FeedEvents []*alerts.Event
    Events []*alerts.Event
    Name string
    EventIndex int
    Markdown string
    AlertId string
    Title string
    Tlp bool
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
	Tlp: true,
    }
}

func (e *State) GetFeedEvent(index int) *alerts.Event {
    if index < 0 {
	index = 0
    } else if index >= len(e.FeedEvents) {
	index = len(e.FeedEvents) - 1
    }

    return e.FeedEvents[index + (e.EventIndex * 10)]
}
