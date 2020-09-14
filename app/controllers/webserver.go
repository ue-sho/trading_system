package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ue-sho/trading_system/config"
)

var templates = template.Must(template.ParseFiles("app/views/google.html"))

// ハンドラ
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "google.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// http://localhost:8080/chart/ にアクセスできるようになる
func StartWebServer() error {
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil) // localhost:8080
}
