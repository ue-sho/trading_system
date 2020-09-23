package main

import (
	"fmt"
	"time"

	"github.com/ue-sho/trading_system/app/models"
)

func main() {

	/* ログの設定 */
	// utils.LoggingSettings(config.Config.LogFile)
	// apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

	/* リアルタイムで情報を取得 */
	// tickerChannel := make(chan bitflyer.Ticker)
	// go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	// for ticker := range tickerChannel {
	// 	fmt.Println(ticker)
	// 	fmt.Println(ticker.GetMidPrice())
	// 	fmt.Println(ticker.DateTime())
	// 	fmt.Println(ticker.TruncateDateTime(time.Second))
	// 	fmt.Println(ticker.TruncateDateTime(time.Minute))
	// 	fmt.Println(ticker.TruncateDateTime(time.Hour))
	// }

	/* オーダー情報 */
	// order := &bitflyer.Order{
	// 	ProductCode:     config.Config.ProductCode,
	// 	ChildOrderType:  "LIMIT", // or MARKET
	// 	Side:            "BUY",   // or SELL
	// 	Price:           7000,
	// 	Size:            0.01,
	// 	MinuteToExpires: 1,
	// 	TimeInForce:     "GTC", // キャンセルするまで有効
	// }

	/* オーダーを出す */
	// res, _ := apiClient.SendOrder(order)
	// fmt.Println(res.ChildOrderAcceptanceID)

	/* オーダー情報をリストをみる */
	// i := "JRF20181012-144016-140584"
	// params := map[string]string{
	//	"product_code": config.Config.ProductCode,
	//	"child_order_acceptance_id": i,
	// }
	// r, _ := apiClient.ListOrder(params)
	// fmt.Println(r)

	// テーブルを出力する
	// fmt.Println(models.DbConnection)

	/* リアルタイムにBitFlyerからデータを取ってきてデータベースに保存する */
	// controllers.StreamIngestionData()
	// controllers.StartWebServer()

	s := models.NewSignalEvents()
	df, _ := models.GetAllCandle("BTC_JPY", time.Minute, 10)
	c1 := df.Candles[0]
	c2 := df.Candles[5]
	s.Buy("BTC_JPY", c1.Time.UTC(), c1.Close, 1.0, true)
	s.Sell("BTC_JPY", c2.Time.UTC(), c2.Close, 1.0, true)
	fmt.Println(models.GetSignalEventsByCount(1))
	fmt.Println(models.GetSignalEventsAfterTime(c1.Time))
	fmt.Println(s.CollectAfter(time.Now().UTC()))
	fmt.Println(s.CollectAfter(c1.Time))
}
