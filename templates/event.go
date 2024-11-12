package templates

import (
	"strconv"

	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/types"
)

func EventList(events []*alerts.Event, index int) string {

	data := struct {
		Events     []*alerts.Event
		EventIndex int
		NextIndex  int
		PrevIndex  int
	}{
		Events:     paginate(events, index),
		EventIndex: index,
		NextIndex:  index + 1,
		PrevIndex:  index - 1,
	}

	const page = `
	<h1>Event: {{.EventIndex}} </h1>
	<button class="outline secondary" id="reset" hx-put="/event/reset" hx-target="body" hx-indicator="#load">Clear Event Cache</button>
	<div id="load" class="htmx-indicator center" aria-busy="true">Loading...</div>
	{{if eq (len .Events) 0}}
	<h2>No New Events</h2>
	{{end}}
	{{range $index, $event := .Events}}
		<article>
			<header>{{$event.Name}}</header>
			{{$event.Ip}}
			<br>
			{{$event.Desc}}
			<br>
			<button class="outline" hx-get="/event/{{$index}}" hx-push-url="true" hx-target="body" onclick="window.scrollTo(0, 0);">Details</button>
		</article>
	{{end}}

	<div class="grid">
		<button id="prev" hx-get="/event/page/{{.PrevIndex}}" hx-push-url="true" hx-target="body" onclick="window.scrollTo(0, 0);"><</button>
		<button id="next" hx-get="/event/page/{{.NextIndex}}" hx-push-url="true" hx-target="body" onclick="window.scrollTo(0, 0);">></button>
		<div></div>
		<div></div>
	</div>
        `

	return ExecuteText("event", page, data)
}

func EventView(event *alerts.Event, index int, form types.Form, eventPage int) string {
	data := struct {
		Name       string
		Event      *alerts.Event
		EventIndex int
		EventPage  int
		Form       string
		FormName   string
	}{
		Name:       event.Name,
		Event:      event,
		EventPage:  eventPage,
		EventIndex: index,
		Form:       getForm(form, event.Name, []*alerts.Event{event}, "/event/"+strconv.Itoa(index)),
		FormName:   types.FormName[form],
	}

	const page = `
        <h1>Event</h1>
	<button hx-get="/event/page/{{.EventPage}}" hx-push-url="true" hx-target="body"><</button>
	<h1>{{.Name}}</h1>
	<h6>{{.Event.Desc}}</h6>
	<article>
		<header>
			<h3>{{.Event.Ip}}</h3>
			<br>
			<small><a href="{{.Event.HostLink}}" target=_blank>Host Link</a></small>
		</header>
		{{if eq (len .Event.Ports) 0}}
			<h4> No Available Information</h4>
		{{end}}
		{{range $key, $value := $.Event.Ports}}
			<h4>{{$key}}</h4>
			{{range $value}}
				<small>{{.Name}}: Priority {{.Rank}}</small>
				<br>
			{{end}}
			<hr>
		{{end}}
	</article>
	<hr>
	<div class="grid">
		<button hx-get="/event/{{.EventIndex}}" hx-target="body">Open Port</button>
		<button hx-get="/event/eol/{{.EventIndex}}" hx-target="body">End of Life</button>
		<button hx-get="/event/login/{{.EventIndex}}" hx-target="body">Login Pages</button>
	</div>
	<h3>{{$.FormName}}</h3>
	{{.Form}}
        `

	return ExecuteText("eventPager", page, data)
}

func paginate(events []*alerts.Event, index int) []*alerts.Event {
	const pageSize int = 10
	entryNum := pageSize * index

	if entryNum >= len(events) {
		return []*alerts.Event{}
	}

	endEntry := min(entryNum+pageSize, len(events))

	return events[entryNum:endEntry]
}
