package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/vikstrous/zengge-lightcontrol/control"
	"github.com/vikstrous/zengge-lightcontrol/local"
)

func setLight(on bool) {
	log.Printf("Light: %v", on)
	host := "192.168.1.103:5577"
	transport, err := local.NewTransport(host)
	if err != nil {
		log.Println("WARN: Failed to connect.")
		return
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

func verifyAndToggleLight() {
	cmd := exec.Command("fuser", "/dev/video0")
	out, err := cmd.Output()
	setLight(err == nil && string(out) != "")
}

func main() {
	// home, err := os.UserHomeDir()
	// panicIf(err)

	verifyAndToggleLight()

	cmd := exec.Command("inotifywait", "/dev/video0", "-q", "-e", "close", "-e", "open", "-m")
	stdout, err := cmd.StdoutPipe()
	panicIf(err)
	err = cmd.Start()
	panicIf(err)

	scanner := bufio.NewScanner(stdout)

	var timer *time.Timer

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		if strings.Contains(line, "OPEN") || strings.Contains(line, "CLOSE") {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(2*time.Second, verifyAndToggleLight)
		}
	}
}
