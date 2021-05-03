package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vikstrous/zengge-lightcontrol/control"
	"github.com/vikstrous/zengge-lightcontrol/local"
)

func setLight(on bool) {
	log.Printf("Light: %v", on)
	host := "192.168.1.103:5577"
	transport, err := local.NewTransport(host)
	if err != nil {
		log.Panicf("Failed to connect. %s", err)
	}

	controller := &control.Controller{transport}
	controller.SetPower(on)
	controller.Close()
}

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setupInterruptHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		setLight(false)
		os.Exit(0)
	}()
}

func main() {
	// home, err := os.UserHomeDir()
	// panicIf(err)

	setupInterruptHandler()

	cmd := exec.Command("inotifywait", "/dev/video0", "-q", "-e", "close", "-e", "open", "-m")
	stdout, err := cmd.StdoutPipe()
	panicIf(err)
	err = cmd.Start()
	panicIf(err)

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		if strings.Contains(line, "OPEN") {
			setLight(true)
		} else if strings.Contains(line, "CLOSE") {
			setLight(false)
		}
	}

	// Hack to allow interrupt handler to exit the application
	time.Sleep(3 * time.Second)
}
