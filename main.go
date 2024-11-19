package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/eagledb14/form-scanner/types"
)

func main() {
	loadEnvVars()
	// fmt.Println(alerts.ParseOtherCreds("1,178347108237rjdfp;laskjdf;lkuj,2,3,4,5"))
	run()
	// o := createform.Osint{
	// 	Name: "Jordan University",
	// 	InScope: []string{"216.237.226.5", "216.237.215.180", "199.184.89.200", "199.184.89.242", "216.237.226.1", "199.184.89.23", "216.237.226.7", "216.237.217.243", "199.184.89.199"},
	// 	OutScope: []string{"199.184.89.239", "216.237.217.244", "216.237.224.6"},
	// 	Events: alerts.DownloadIpList("Jordan University", "199.184.89.242, 199.184.89.242, 199.184.89.235,  216.237.226.26"),
	// }
	// file, _ := os.ReadFile("cred-test")
	// creds := alerts.ParseCredentialDump(string(file), "nhcgov.com")
	// o.Creds = creds
	//
	// html := createform.CreateCoverHtml(o.CreateMarkdown(), o.Name)
	// _ = html
	// fmt.Println(html)

}

func run() {
	checkResources()

	auto := flag.Bool("auto", false, "run in automatic mode")
	flag.Parse()

	if *auto {
		autoCreateEventFiles()
	} else {
		state := types.NewState()
		var port = ""

if os.Getenv("DEV") == "true" {
			port = ":8080"
		} else {
			port, _ = getRandomPort()
		}

		go openBrowser("localhost" + port)
		serv(port, &state)
	}
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

func getRandomPort() (string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return ":" + strconv.Itoa(addr.Port), nil
}

func checkResources() {
	path := "./resources"

	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		panic("Missing resources folder")
	}
}
