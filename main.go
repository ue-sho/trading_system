package main

import (
	"fmt"

	"github.com/ue-sho/trading_system/app/models"
	"github.com/ue-sho/trading_system/config"
	"github.com/ue-sho/trading_system/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	//apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

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
	//i := "JRF20181012-144016-140584"
	//params := map[string]string{
	//	"product_code": config.Config.ProductCode,
	//	"child_order_acceptance_id": i,
	//}
	//r, _ := apiClient.ListOrder(params)
	//fmt.Println(r)

	// テーブルを出力する
	fmt.Println(models.DbConnection)
}
