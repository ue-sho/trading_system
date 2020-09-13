package main

import (
	"fmt"

	"github.com/ue-sho/trading_system/bitflyer"
	"github.com/ue-sho/trading_system/config"
	"github.com/ue-sho/trading_system/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	fmt.Println(apiClient.GetBalance())
}
