package janitor

import (
	"time"

	"github.com/Sirupsen/logrus"
)

func DefaultConfig() Config {
	config := Config{
		ListenAddr:    "0.0.0.0:80",
		FlushInterval: time.Second * 1,
		Domain:        "lvh.me",
		LogLevel:      logrus.InfoLevel,
	}

	return config
}

type Config struct {
	ListenAddr    string
	FlushInterval time.Duration
	Domain        string
	LogLevel      logrus.Level
}
