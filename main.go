package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func SetupLogger() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	janitorConfig := DefaultConfig()
	//enable multi_port mode
	//janitorConfig.Listener.Mode = config.MULTIPORT_LISTENER_MODE

	//TuneGolangProcess()
	SetupLogger()

	server := NewJanitorServer(janitorConfig)
	go server.ServerInit().Run()

	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		log.Debug("sending targetChangeEvent")
		time.Sleep(time.Second * 10)
		targetChangeEvents := []*TargetChangeEvent{
			{
				Change:   "add",
				App:      "nginx0051",
				AppID:    "nginx0051-xcm-datamanmesos",
				TaskID:   "0-nginx0051-xcm-datamanmesos",
				TaskIp:   "192.168.1.162",
				PortName: "web",
				TaskPort: 80,
			},
			{
				Change:   "add",
				App:      "nginx0051",
				AppID:    "nginx0051-xcm-datamanmesos",
				TaskID:   "1-nginx0051-xcm-datamanmesos",
				TaskIp:   "192.168.1.162",
				PortName: "web1",
				TaskPort: 80,
			},
		}

		for _, targetChangeEvent := range targetChangeEvents {
			server.SwanEventChan() <- targetChangeEvent
		}
		time.Sleep(time.Second * 10)
		targetChangeEvents = []*TargetChangeEvent{
			{
				Change:   "del",
				App:      "nginx0051",
				AppID:    "nginx0051-xcm-datamanmesos",
				TaskID:   "0-nginx0051-xcm-datamanmesos",
				TaskIp:   "192.168.1.162",
				PortName: "web",
				TaskPort: 80,
			},
			{
				Change:   "del",
				App:      "nginx0051",
				AppID:    "nginx0051-xcm-datamanmesos",
				TaskID:   "1-nginx0051-xcm-datamanmesos",
				TaskIp:   "192.168.1.162",
				PortName: "web1",
				TaskPort: 80,
			},
		}
		for _, targetChangeEvent := range targetChangeEvents {
			server.SwanEventChan() <- targetChangeEvent
		}
	}
}
