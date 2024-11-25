package alerts

import (
	"time"
)

type Feed struct {
	events []*Event
	Index int
}

func NewFeed() Feed {
	events := DownloadRss()
	return Feed{
		events: events,
		Index: 0,
	}
}

func (f *Feed) Next() *Event {
	f.Index = min(f.Index + 1, len(f.events) - 1)
	return f.GetEvent()
}

func (f *Feed) GetEvent() *Event {
	for f.events[f.Index].Loaded != true {
		time.Sleep(1)
	}

	return f.events[f.Index]
}

func (f *Feed) Prev() *Event {
	f.Index = max(f.Index - 1, 0)
	return f.GetEvent()
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
