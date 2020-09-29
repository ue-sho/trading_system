package controllers

import (
	"log"

	"github.com/ue-sho/trading_system/app/models"
	"github.com/ue-sho/trading_system/bitflyer"
	"github.com/ue-sho/trading_system/config"
)

// ストリーミングでデータを取ってくる
func StreamIngestionData() {
	c := config.Config // 名前が長いので宣言
	ai := NewAI(c.ProductCode, c.TradeDuration, c.DataLimit, c.UsePercent, c.StopLimitPercent, c.BackTest)

	var tickerChannl = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannl)

	go func() { // 違うスレッドでfor文を回すことで他のことをできるようにする
		for ticker := range tickerChannl {
			log.Printf("action=StreamIngestionData, %v", ticker)
			for _, duration := range config.Config.Durations { // 1秒、1分、１時間のものそれぞれデータベースのテーブルに書き込む
				isCreated := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreated == true && duration == config.Config.TradeDuration {
					ai.Trade()
				}
			}
		}
	}() //この丸括弧がないと定義しただけで実行できない
}
