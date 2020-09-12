package main

import (
	"fmt"

	"github.com/ue-sho/trading_system/config"
)

func main() {
	fmt.Println(config.Config.ApiKey)
	fmt.Println(config.Config.ApiSecret)
}
