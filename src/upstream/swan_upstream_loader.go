package upstream

import (
	"fmt"
	"net"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

const (
	SWAN_UPSTREAM_LOADER_KEY = "SwanUpstreamLoader"
)

type AppEventNotify struct {
	Operation     string
	TaskName      string
	AgentHostName string
	AgentPort     string
}

type SwanUpstreamLoader struct {
	UpstreamLoader
	Upstreams    []*Upstream
	changeNotify chan bool
	sync.Mutex
	swanEventChan     chan *AppEventNotify
	DefaultUpstreamIp net.IP
	Port              string
	Proto             string
}

func InitSwanUpstreamLoader(defaultUpstreamIp net.IP, defaultPort string) (*SwanUpstreamLoader, error) {
	swanUpstreamLoader := &SwanUpstreamLoader{}
	swanUpstreamLoader.changeNotify = make(chan bool, 64)
	swanUpstreamLoader.Upstreams = make([]*Upstream, 0)
	swanUpstreamLoader.DefaultUpstreamIp = defaultUpstreamIp
	swanUpstreamLoader.Port = defaultPort
	swanUpstreamLoader.Proto = "http"
	swanUpstreamLoader.swanEventChan = make(chan *AppEventNotify, 1)
	go swanUpstreamLoader.Poll()
	return swanUpstreamLoader, nil
}

func (swanUpstreamLoader *SwanUpstreamLoader) Poll() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("SwanUpstreamLoader poll got error: %s", err)
			swanUpstreamLoader.Poll() // execute poll again
		}
	}()

	for {
		//var appEvent *AppEventNotify
		log.Debug("upstreamLoader is listening app event...")
		appEvent := <-swanUpstreamLoader.swanEventChan
		log.Debug("upstreamLoader receive one app event:%s", appEvent)
		switch strings.ToLower(appEvent.Operation) {
		case "add":
			newUpstream := buildSwanUpstream(appEvent, swanUpstreamLoader.DefaultUpstreamIp, swanUpstreamLoader.Port, swanUpstreamLoader.Proto)
			fmt.Printf("upstreamKey:%s\n", newUpstream.Key())
			//for _, upstream := range swanUpstreamLoader.Upstreams {
			//}
		case "delete":
		}

	}
}

func (swanUpstreamLoader *SwanUpstreamLoader) List() []*Upstream {
	swanUpstreamLoader.Lock()
	defer swanUpstreamLoader.Unlock()
	return swanUpstreamLoader.Upstreams
}

func (swanUpstreamLoader *SwanUpstreamLoader) SwanEventChan() chan<- *AppEventNotify {
	return swanUpstreamLoader.swanEventChan
}

func (swanUpstreamLoader *SwanUpstreamLoader) ServiceEntries() []string {
	entryList := make([]string, 0)
	for _, u := range swanUpstreamLoader.Upstreams {
		entry := fmt.Sprintf("%s://%s:%s", u.Key().Proto, u.Key().Ip, u.Key().Port)
		entryList = append(entryList, entry)
	}

	return entryList
}

func (swanUpstreamLoader *SwanUpstreamLoader) Get(serviceName string) *Upstream {
	return nil
}

func (swanUpstreamLoader *SwanUpstreamLoader) Remove(upstream *Upstream) {
	index := -1
	for k, v := range swanUpstreamLoader.Upstreams {
		if v == upstream {
			index = k
			break
		}
	}

	if index >= 0 {
		swanUpstreamLoader.Upstreams = append(swanUpstreamLoader.Upstreams[:index], swanUpstreamLoader.Upstreams[index+1:]...)
	}
}

func (swanUpstreamLoader *SwanUpstreamLoader) ChangeNotify() <-chan bool {
	return swanUpstreamLoader.changeNotify
}

func buildSwanUpstream(appEvent *AppEventNotify, defaultUpstreamIp net.IP, port string, proto string) Upstream {
	// create a new upstream
	var upstream Upstream
	taskNamespaces := strings.Split(appEvent.TaskName, ".")
	taskNum := taskNamespaces[0]
	appName := strings.Join(taskNamespaces[1:], ".")
	upstream.ServiceName = appName
	upstream.FrontendIp = defaultUpstreamIp.String()
	upstream.FrontendPort = port
	upstream.FrontendProto = proto
	fmt.Printf("taskNum:%s\n", taskNum)
	fmt.Printf("appName:%s\n", appName)
	upstream.Targets = make([]*Target, 0)
	upstream.StaleMark = false
	upstream.SetState(STATE_NEW)

	// create a new target
	var target Target
	target.Address = appEvent.AgentHostName
	target.ServiceName = appName
	target.ServiceID = taskNum
	target.ServiceAddress = appEvent.AgentHostName
	target.ServicePort = appEvent.AgentPort
	target.Upstream = &upstream

	// add target in upstream
	upstream.Targets = append(upstream.Targets, &target)
	return upstream
}
