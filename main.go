package main

import (
	"time"

	"github.com/op/go-logging"
)

//Config global var to share configuration data
var Config = LoadConfig().(*Configuration)

//AbattoirLog - Logging
var AbattoirLog = logging.MustGetLogger("Abattoir")

func main() {
	for {
		seekAndDestroyBadHosts(Config)
		time.Sleep(time.Duration(Config.RunInterval) * time.Second)
	}
}
