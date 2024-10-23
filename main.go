package main

import (
	"bufio"
	"fmt"

	// "github.com/eagledb14/form-scanner/alerts"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/eagledb14/form-scanner/alerts"
	createform "github.com/eagledb14/form-scanner/create-form"
	t "github.com/eagledb14/form-scanner/templates"
	"github.com/eagledb14/form-scanner/types"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// fmt.Println("Hellow Worl")
	loadEnvVars()
	// fmt.Println(alerts.Download())
	// ch := make(chan int)
	// for _, a := range alerts.Download() {
	// 	// a.GetName(0)
	// 	// a.GetAlertId(0)
	// 	// a.GetName(0)
	// 	// fmt.Println(a.Name)
	//
	// 	go a.Load()
	// 	// fmt.Println(a)
	// 	// beak
	// }
	// events := alerts.Download()
	// fmt.Println(len(events))
	// <-ch
	// fmt.Println(os.Getenv("API_KEY"))
	// alerts.Download()
	// events := alerts.DownloadIpList("monkey", "24.246.129.0/24")
	// fmt.Println(len(events))
	// for _, e := range events {
	// 	fmt.Println(len(e.Ports))
	// }
	// fmt.Println(events)
	state := types.NewState()

	go openBrowser("localhost:8080")
	serv(":8080", &state)
}

func serv(port string, state *types.State) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak()))
	})

	app.Get("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak()))
	})

	app.Post("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		form := createform.CredLeak{
			OrgName:    c.FormValue("orgName"),
			FormNumber: c.FormValue("formNumber"),
			VictimOrg:  c.FormValue("victimOrg"),
			Leaks:      c.FormValue("leaks"),
			Creds:      c.FormValue("creds"),
			Password:   c.FormValue("password"),
			IpAddress:  c.FormValue("ipAddress"),
			UserPass:   c.FormValue("userPass"),
			AddInfo:    c.FormValue("addInfo"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp"),
		}
		_ = form

		return c.SendString(t.BuildPage(t.Index()))
	})

	app.Get("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		if len(state.Events) > 0 {
			return c.SendString(t.BuildPage(t.OpenPortForm(state.Name, state.Events)))
		}

		return c.SendString(t.BuildPage(t.OpenPortDownload()))
	})

	//makes a new file
	app.Post("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		fmt.Println("yay")
		form := createform.OpenPort {
			FormNumber: c.FormValue("formNumber"),
			Threat: c.FormValue("threat"),
			Summary: c.FormValue("summary"),
			Body: c.FormValue("body"),
			Reference: c.FormValue("reference"),
			Tlp: c.FormValue("tlp"),
			Events: state.Events,
		}
		_ = form

		return c.SendString(t.BuildPage("yay, you made the page, congrats"))
	})

	// clears the state of the thing
	app.Put("/openport", func(c *fiber.Ctx) error {
		state.Events = []*alerts.Event{}
		state.Name = ""
		return c.SendString(t.BuildPage(t.OpenPortDownload()))
	})

	app.Post("/openportform", func(c *fiber.Ctx) error {
		name := c.FormValue("orgName")
		ips := c.FormValue("ipAddress")

		events := alerts.DownloadIpList(name, ips)
		events = alerts.FilterEvents(events)

		state.Events = events
		state.Name = name

		return c.SendString(t.BuildPage(t.OpenPortForm(state.Name, state.Events)))
	})

	app.Get("/actor", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.Actors()))
	})

	app.Get("/event", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.Event()))
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

func loadEnvVars() error {
	file, err := os.Open("./resources/key.env")
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
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
			return fmt.Errorf("invalid line format: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting env var %s: %v", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}
