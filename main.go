package main

import (
	"log"

	"github.com/ue-sho/trading_system/config"
	"github.com/ue-sho/trading_system/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("test")
}
