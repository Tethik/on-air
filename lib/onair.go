package lib

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vikstrous/zengge-lightcontrol/control"
	"github.com/vikstrous/zengge-lightcontrol/local"
)

func SetLight(on bool) {
	log.Printf("Light: %v", on)
	host := viper.GetString("host")

	if len(host) == 0 {
		log.Println("WARN: host is undefined. Have you created the configuration file? ", viper.ConfigFileUsed())
		return
	}

	transport, err := local.NewTransport(host)
	if err != nil {
		log.Printf("WARN: Failed to connect to %s\n", host)
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

func VerifyAndToggleLight() {
	device := viper.GetString("device")
	if len(device) == 0 {
		log.Println("WARN: device is undefined. Have you created the configuration file?")
		return
	}
	cmd := exec.Command("fuser", device)
	out, err := cmd.Output()
	SetLight(err == nil && string(out) != "")
}

func Daemon() {
	// home, err := os.UserHomeDir()
	// panicIf(err)

	VerifyAndToggleLight()

	device := viper.GetString("device")
	cmd := exec.Command("inotifywait", device, "-q", "-e", "close", "-e", "open", "-m")
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
			timer = time.AfterFunc(2*time.Second, VerifyAndToggleLight)
		}
	}
}
