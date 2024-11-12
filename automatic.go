package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eagledb14/form-scanner/alerts"
	createform "github.com/eagledb14/form-scanner/create-form"
	"github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
)

func autoCreateEventFiles() {
	fmt.Println("Generating...")
	os.MkdirAll("generated-forms", 0755)

	events := alerts.DownloadRss()
	for _, e := range events {
		time.Sleep(time.Duration(3 * time.Second))
		go func(e *alerts.Event) {
			e.Load()
		}(e)
	}
	forms := []createform.OpenPort{}

	for i, e := range events {
		form := createform.OpenPort{
			OrgName:    e.Name,
			FormNumber: strconv.Itoa(i),
			Threat:     "T1133 External Remote Services",
			Summary:    templates.OpenPortSummary(e.Name, []*alerts.Event{e}),
			Body:       templates.OpenPortBody(e.Name, []*alerts.Event{e}),
			Tlp:        true,
			Events:     []*alerts.Event{e},
		}
		forms = append(forms, form)
	}

	state := &types.State{}
	for i, form := range forms {
		md := form.CreateMarkdown(state)
		html := createform.CreateHtml(md, events[i].Name, true)

		fileName := "./generated-forms/" + events[i].Name + "-" + state.AlertId + ".html"
		fmt.Println(fileName)

		file, _ := os.Create(fileName)
		file.WriteString(html)
		file.Close()
	}
}
