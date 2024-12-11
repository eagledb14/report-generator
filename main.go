package main

import (
	"bufio"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/eagledb14/form-scanner/types"
)

func main() {
	loadEnvVars()
	run()
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

		serv(port, state)
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
