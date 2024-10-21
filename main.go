package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("Hellow Worl")
	// go openBrowser("localhost:8080")
	serv(":8080")
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
