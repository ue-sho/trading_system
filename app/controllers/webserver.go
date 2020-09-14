package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ue-sho/trading_system/app/models"
	"github.com/ue-sho/trading_system/config"
)

var templates = template.Must(template.ParseFiles("app/views/google.html"))

// ハンドラ
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	limit := 100
	duration := "1m"
	durationTime := config.Config.Durations[duration] // mapなので time.Minute を取ってくる
	df, err := models.GetAllCandle(config.Config.ProductCode, durationTime, limit)
	if err != nil {
		log.Fatalf("GetAllCandle Error : %s", err)
	}

	err = templates.ExecuteTemplate(w, "google.html", df.Candles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// http://localhost:8080/chart/ にアクセスできるようになる
func StartWebServer() error {
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil) // localhost:8080
}
