package main

import (
	"bufio"
	"fmt"

	// "sync"

	// "github.com/eagledb14/form-scanner/alerts"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/eagledb14/form-scanner/alerts"
	createform "github.com/eagledb14/form-scanner/create-form"
	t "github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
	"github.com/gofiber/fiber/v2"
)

func main() {
	loadEnvVars()
	// events := alerts.DownloadRss()
	// wg := sync.WaitGroup{}
	// for _, e := range events {
	// 	wg.Add(1)
	// 	go func(e *alerts.Event, wg *sync.WaitGroup) {
	// 		e.Load()
	// 		wg.Done()
	// 	}(e, &wg)
	// }
	// wg.Wait()
	//
	// form := createform.OpenPort{
	// 	OrgName: "Test Title",
	// 	FormNumber: "1",
	// 	Threat: "2",
	// 	Summary: "hias fsealksjdfa;sdjf askdfj",
	// 	Body: "aksjdflkajsdf alksjdf",
	// 	Tlp: false,
	// 	Reference: "",
	// 	// Events: alerts.DownloadIpList("Testing thing", "24.172.113.143"),
	// 	Events: events, 
	// }
	// 
	// fmd := form.CreateMarkdown()
	// createform.CreateHtml(fmd, "Threat Intel Summary")
	// createform.CreateHtml(createform.Header("threat summary", true))
	
	
	state := types.NewState()

	go openBrowser("localhost:8080")
	serv(":8080", &state)
}

func serv(port string, state *types.State) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak(), state))
	})

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

		return c.Redirect("/preview")
	})

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
			OrgName: state.Name,
			FormNumber: c.FormValue("formNumber"),
			Threat:     c.FormValue("threat"),
			Summary:    c.FormValue("summary"),
			Body:       c.FormValue("body"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
			Events:     state.Events,
		}
		state.Markdown = form.CreateMarkdown(state)

		return c.Redirect("/preview")
	})

	// clears the state of the thing
	app.Put("/openport", func(c *fiber.Ctx) error {
		state.Events = []*alerts.Event{}
		state.Name = ""
		return c.SendString(t.BuildPage(t.OpenPortDownload(), state))
	})

	app.Post("/openport/form", func(c *fiber.Ctx) error {
		name := c.FormValue("orgName")
		ips := c.FormValue("ipAddress")

		events := alerts.DownloadIpList(name, ips)
		events = alerts.FilterEvents(events)

		state.Events = events
		state.Name = strings.Clone(name)
		// state.SetName(strings.Clone(name), "post openport/form")

		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/port", func(c *fiber.Ctx) error {
		// state.FormType = types.Open
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/eol", func(c *fiber.Ctx) error {
		// state.FormType = types.EOL
		return c.SendString(t.BuildPage(t.OpenPortForm(types.EOL, state.Name, state.Events), state))
	})

	app.Get("/openport/login", func(c *fiber.Ctx) error {
		// state.FormType = types.Login
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Login, state.Name, state.Events), state))
	})

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
		state.Name = strings.Clone(form.Name)
		// state.SetName(strings.Clone(form.Name), "post /actor")

		return c.Redirect("/preview")
	})

	app.Get("/event/page/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 {
			index = 0
		}

		if index >= len(state.FeedEvents)/10 {
			index = ((len(state.FeedEvents) / 10) / 2) * 2
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
		// state.FormType = types.Open

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
		// state.FormType = types.Open

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
			OrgName: state.FeedEvents[index].Name,
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

		return c.Redirect("/preview")
	})

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
		c.Set("Content-Disposition", "attachment; filename=\"" + state.Name + "-" + state.AlertId +".html\"")

		form := createform.CreateHtml(state.Markdown, "Threat Intel Summary")
		return c.SendString(form)
	})

	app.Static("/style.css", "./resources/style.css")

	app.Listen(port)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Error opening browser:", err)
	}
}

func loadEnvVars() {
	file, err := os.Open("./resources/key.env")
	if err != nil {
		panic("error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			panic("invalid line format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := os.Setenv(key, value); err != nil {
			panic("error setting env var")
		}
	}

	if err := scanner.Err(); err != nil {
		panic("error reading file")
	}
}
