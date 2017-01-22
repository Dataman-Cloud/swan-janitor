package main

import (
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

type Target struct {
	App      string
	AppID    string
	TaskID   string
	TaskIp   string
	TaskPort uint32
	PortName string
	Upstream *Upstream
}

func (t *Target) Equal(t1 *Target) bool {
	return t.App == t1.App &&
		t.AppID == t1.AppID &&
		t.TaskID == t1.TaskID &&
		t.TaskIp == t1.TaskIp &&
		t.TaskPort == t1.TaskPort &&
		t.PortName == t1.PortName
}

func (t *Target) ToString() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", t.App,
		t.AppID, t.TaskID, t.TaskIp, t.TaskPort, t.PortName)
}

func (t Target) Entry() *url.URL {
	taskEntry := fmt.Sprintf("http://%s:%d", t.TaskIp, t.TaskPort)
	url, err := url.Parse(taskEntry)
	if err != nil {
		log.Error("parse target entry %s to url got err %s", taskEntry, err)
	}

	return url
}
