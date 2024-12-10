package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/eagledb14/form-scanner/alerts"
	createform "github.com/eagledb14/form-scanner/create-form"
	t "github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
	"github.com/gofiber/fiber/v2"
)

func serv(port string, state *types.State) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak(), state))
	})

	servCredLeak(app, state)
	servOpenPort(app, state)
	servActor(app, state)
	servEvents(app, state)
	servEvents(app, state)
	servMarkdown(app, state)
	servCsv(app, state)
	servPortViewer(app, state)
	servOsint(app, state)

	app.Static("/style.css", "./resources/style.css")

	app.Listen(port)
}

func servCredLeak(app *fiber.App, state *types.State) {
	app.Get("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak(), state))
	})

	app.Post("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		form := createform.CredLeak{
			OrgName:    c.FormValue("orgName"),
			FormNumber: c.FormValue("formNumber"),
			VictimOrg:  c.FormValue("victimOrg"),
			Password:   c.FormValue("password"),
			UserPass:   c.FormValue("userPass"),
			AddInfo:    c.FormValue("addInfo"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.OrgName)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp
		state.Report = types.Header

		return c.Redirect("/preview")
	})
}

func servOpenPort(app *fiber.App, state *types.State) {
	app.Get("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		if len(state.Events) > 0 {
			return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
		}

		return c.SendString(t.BuildPage(t.OpenPortDownload(), state))
	})

	//makes a new file
	app.Post("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		form := createform.OpenPort{
			OrgName:    state.Name,
			FormNumber: c.FormValue("formNumber"),
			Threat:     c.FormValue("threat"),
			Summary:    c.FormValue("summary"),
			Body:       c.FormValue("body"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
			Events:     state.Events,
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp
		state.Report = types.Header

		return c.Redirect("/preview")
	})

	// clears the state of the selected event
	app.Put("/openport", func(c *fiber.Ctx) error {
		state.Events = []*alerts.Event{}
		state.Name = ""
		return c.SendString(t.BuildPage(t.OpenPortDownload(), state))
	})

	app.Post("/openport/form", func(c *fiber.Ctx) error {
		name := c.FormValue("orgName")
		ips := c.FormValue("ipAddress")

		events := alerts.DownloadIpList(name, ips)

		state.Events = events
		state.Name = strings.Clone(name)

		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/port", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/eol", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.EOL, state.Name, state.Events), state))
	})

	app.Get("/openport/login", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Login, state.Name, state.Events), state))
	})
}

func servActor(app *fiber.App, state *types.State) {
	app.Get("/actor", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.Actors(), state))
	})

	app.Post("/actor", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		form := createform.Actor{
			Name:         c.FormValue("name"),
			Alias:        c.FormValue("alias"),
			Date:         c.FormValue("date"),
			Country:      c.FormValue("country"),
			Motivation:   c.FormValue("motivation"),
			Target:       c.FormValue("target"),
			Malware:      c.FormValue("malware"),
			Reporter:     c.FormValue("report"),
			Confidence:   c.FormValue("confidence"),
			Exploits:     c.FormValue("exploits"),
			Summary:      c.FormValue("summary"),
			Capabilities: c.FormValue("capabilities"),
			Detection:    c.FormValue("detection"),
			Ttps:         c.FormValue("ttps"),
			Infra:        c.FormValue("infra"),
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.Name)
		state.Title = "Threat Actor Profile"
		state.Tlp = false
		state.Report = types.Header

		return c.Redirect("/preview")
	})
}

func servEvents(app *fiber.App, state *types.State) {
	app.Get("/event/page/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 {
			index = 0
		}

		state.EventIndex = index

		return c.SendString(t.BuildPage(t.EventList(state.FeedEvents, index), state))
	})

	app.Get("/event/open/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Open, state.EventIndex), state))
	})

	app.Get("/event/eol/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.EOL, state.EventIndex), state))
	})

	app.Get("/event/login/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Login, state.EventIndex), state))
	})

	app.Get("/event/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Open, state.EventIndex), state))
	})

	app.Post("/event/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		form := createform.OpenPort{
			OrgName:    state.FeedEvents[index].Name,
			FormNumber: c.FormValue("formNumber"),
			Threat:     c.FormValue("threat"),
			Summary:    c.FormValue("summary"),
			Body:       c.FormValue("body"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
			Events:     []*alerts.Event{state.FeedEvents[index]},
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.OrgName)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp
		state.Report = types.Header

		return c.Redirect("/preview")
	})

	app.Put("/event/reset", func(c *fiber.Ctx) error {
		cache := alerts.NewEventCache()
		cache.ClearTable()

		events := alerts.DownloadRss()
		state.FeedEvents = events
		state.LoadEvents()
		state.EventIndex = 0
		time.Sleep(time.Duration(2 * time.Second))

		return c.SendString(t.BuildPage(t.EventList(state.FeedEvents, state.EventIndex), state))
	})
}

func servMarkdown(app *fiber.App, state *types.State) {
	app.Get("/preview", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return c.SendString(t.BuildPage(t.MarkdownViewer(state), state))
	})

	app.Post("/preview", func(c *fiber.Ctx) error {
		md := c.FormValue("markdown")
		state.Markdown = md

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/create", func(c *fiber.Ctx) error {
		c.Set("Content-Disposition", "attachment; filename=\""+state.Name+"-"+state.AlertId+".html\"")

		form := ""

		switch state.Report {
		case types.Header:
			form = createform.CreateHeaderHtml(state.Markdown, state.Title, state.Tlp)
		case types.Cover:
			form = createform.CreateCoverHtml(state.Markdown, state.Title)
		}

		return c.SendString(form)
	})

}

func servCsv(app *fiber.App, state *types.State) {
	app.Get("/csv", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return c.SendString(t.BuildPage(t.Csv(), state))
	})

	app.Post("/csv", func(c *fiber.Ctx) error {
		name := c.FormValue("orgName")
		query := c.FormValue("ipAddress")

		state.Name = strings.Clone(name)
		state.Markdown = createform.CreateCsv(query)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/csv/create", func(c *fiber.Ctx) error {
		c.Set("Content-Disposition", "attachment; filename=\""+state.Name + ".csv\"")
		return c.SendString(state.Markdown)
	})
}

func servPortViewer(app *fiber.App, state *types.State) {
	app.Get("/portview", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.PortViewer(), state))
	})

	app.Post("/portview", func(c *fiber.Ctx) error {
		ips := c.FormValue("ipAddress")
		form := createform.PortViewer{
			Events: alerts.DownloadIpList("", ips),
		}

		state.Markdown = form.CreateMarkdown()
		state.Name = ""
		state.Tlp = false

		return c.Redirect("/preview")
	})
}

func servOsint(app *fiber.App, state *types.State) {
	app.Get("/osint", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.Osint(), state))
	})

	app.Post("/osint", func(c *fiber.Ctx) error {
		name := strings.Clone(c.FormValue("orgName"))
		inScope := c.FormValue("inScope")
		inScopeList := strings.FieldsFunc(inScope, func(r rune) bool {
			return r == ',' || r == ' '
		})

		outScope := c.FormValue("outScope")
		outScopeList := strings.FieldsFunc(outScope, func(r rune) bool {
			return r == ',' || r == ' '
		})

		inScopEvents := []*alerts.Event{}
		if inScope != "" {
			inScopEvents = alerts.DownloadIpList(name, inScope)
		}

		outScopeEvents := []*alerts.Event{}
		if outScope != "" {
			outScopeEvents = alerts.DownloadIpList(name, outScope)
		}

		events := append(outScopeEvents, inScopEvents...)
		events = alerts.FilterCveEvents(events)

		recordedFutureCreds := alerts.ParseCredentialDump(c.FormValue("recordedFutureCreds"))
		otherCreds := alerts.ParseOtherCreds(c.FormValue("otherCreds"))

		vulnerableUrls, _ := strconv.Atoi(c.FormValue("vulnerableUrls"))

		creds := append(recordedFutureCreds, otherCreds...)
		creds = alerts.SortCreds(creds)

		form := createform.Osint{
			Name: name, 
			InScope: inScopeList,
			OutScope: outScopeList,
			Events: events,
			Creds: creds,
			Url: c.FormValue("url"),
			VulnerableUrls: vulnerableUrls,
			AssetSeverity: c.FormValue("assetSeverity"),
			AccountSeverity: c.FormValue("accountSeverity"),
			WebsiteSeverity: c.FormValue("websiteSeverity"),
		}
		state.Name = form.Name
		state.Title = ""
		state.Report = types.Cover
		state.Markdown = form.CreateMarkdown()


		return c.Redirect("/preview")
	})
}
