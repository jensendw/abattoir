package main

import (
	"github.com/op/go-logging"
  "time"
)

//Config global var to share configuration data
var Config = LoadConfig().(*Configuration)

//Logging
var log = logging.MustGetLogger("Abattoir")

func main() {
  for {
    seekAndDestroyBadHosts()
    time.Sleep(time.Duration(Config.RunInterval) * time.Second)
  }
}
