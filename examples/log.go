package main

import (
	"time"

	log "github.com/zzh20/logger"
)

func main() {
	log.Init("examples/logger.json")
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
	time.Sleep(time.Second * 2)
}
