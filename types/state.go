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

    newState :=  State{
	FeedEvents: feedEvents,
	Tlp: true,
    }

    newState.LoadEvents()
    return newState
}

func (e *State) LoadEvents() {
    for _, e := range e.FeedEvents {
	time.Sleep(time.Duration(1 * time.Second))
	go func(e *alerts.Event) {
	    e.Load()
	}(e)
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
