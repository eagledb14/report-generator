package main

import (
	"bufio"
	"fmt"
	"github.com/eagledb14/form-scanner/alerts"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	// fmt.Println("Hellow Worl")
	LoadEnvVars()
	// fmt.Println(alerts.Download())
	ch := make(chan int)
	for _, a := range alerts.Download() {
		// a.GetName(0)
		// a.GetAlertId(0)
		// a.GetName(0)
		// fmt.Println(a.Name)

		go a.Load()
		// fmt.Println(a)
		// beak
	}
	<-ch
	// fmt.Println(os.Getenv("API_KEY"))
	// alerts.Download()

	// go openBrowser("localhost:8080")
	// serv(":8080")
}

func serv(port string) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("hi")
	})

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

func LoadEnvVars() error {
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
