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
    Report ReportType
}

func NewState() *State {

    newState := &State{
	Tlp: true,
	Report: Header,
    }

    go func(state *State) {
	feedEvents := alerts.DownloadRss()
	newState.FeedEvents = feedEvents
	newState.LoadEvents()
    }(newState)
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

    pagedIndex := index + (e.EventIndex * 10)
    if pagedIndex >= len(e.FeedEvents) {
	pagedIndex = len(e.FeedEvents) - 1
    }

    return e.FeedEvents[pagedIndex]
}
