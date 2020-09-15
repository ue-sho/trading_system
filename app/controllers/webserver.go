package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ue-sho/trading_system/app/models"
	"github.com/ue-sho/trading_system/config"
)

// キャッシュして保存しておく
var templates = template.Must(template.ParseFiles("app/views/chart.html"))

// ハンドラ
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	/* ajaxでデータを取ってくるようにしたのでいらなくなった */
	// limit := 100
	// duration := "1m"                                  // ここで間隔を変えていく
	// durationTime := config.Config.Durations[duration] // mapなので time.Minute を取ってくる
	// df, err := models.GetAllCandle(config.Config.ProductCode, durationTime, limit)
	// if err != nil {
	// 	log.Fatalf("GetAllCandle Error : %s", err)
	// }

	// err = templates.ExecuteTemplate(w, "google.html", df.Candles)
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// json型でエラーを返す
// code : エラーコード 404とか
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code}) // json型に変換
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError) // httpリクエストにjsonを書き込む
}

var apiValidPath = regexp.MustCompile("^/api/candle/$") // 正規表現で表される左の文字列を登録

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := apiValidPath.FindStringSubmatch(r.URL.Path) // 正規表現と一致する部分だけをとる
		if len(m) == 0 {
			APIError(w, "Not found", http.StatusNotFound)
		}
		fn(w, r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		APIError(w, "No product_code param", http.StatusBadRequest)
		return
	}
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Config.Durations[duration]

	df, _ := models.GetAllCandle(productCode, durationTime, limit)

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// http://localhost:8080/chart/ にアクセスできるようになる
func StartWebServer() error {
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil) // localhost:8080
}
