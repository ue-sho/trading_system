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

	sma := r.URL.Query().Get("sma")
	if sma != "" {
		strSmaPeriod1 := r.URL.Query().Get("smaPeriod1")
		strSmaPeriod2 := r.URL.Query().Get("smaPeriod2")
		strSmaPeriod3 := r.URL.Query().Get("smaPeriod3")
		period1, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 7
		}
		period2, err := strconv.Atoi(strSmaPeriod2)
		if strSmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 14
		}
		period3, err := strconv.Atoi(strSmaPeriod3)
		if strSmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 50
		}
		df.AddSma(period1)
		df.AddSma(period2)
		df.AddSma(period3)
	}

	ema := r.URL.Query().Get("ema")
	if ema != "" {
		strEmaPeriod1 := r.URL.Query().Get("emaPeriod1")
		strEmaPeriod2 := r.URL.Query().Get("emaPeriod2")
		strEmaPeriod3 := r.URL.Query().Get("emaPeriod3")
		period1, err := strconv.Atoi(strEmaPeriod1)
		if strEmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 7
		}
		period2, err := strconv.Atoi(strEmaPeriod2)
		if strEmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 14
		}
		period3, err := strconv.Atoi(strEmaPeriod3)
		if strEmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 50
		}
		df.AddEma(period1)
		df.AddEma(period2)
		df.AddEma(period3)
	}

	bbands := r.URL.Query().Get("bbands")
	if bbands != "" {
		strN := r.URL.Query().Get("bbandsN")
		strK := r.URL.Query().Get("bbandsK")
		n, err := strconv.Atoi(strN)
		if strN == "" || err != nil || n < 0 {
			n = 20
		}
		k, err := strconv.Atoi(strK)
		if strK == "" || err != nil || k < 0 {
			k = 2
		}
		df.AddBBands(n, float64(k))
	}

	ichimoku := r.URL.Query().Get("ichimoku")
	if ichimoku != "" {
		df.AddIchimoku()
	}

	rsi := r.URL.Query().Get("rsi")
	if rsi != "" {
		strPeriod := r.URL.Query().Get("rsiPeriod")
		period, err := strconv.Atoi(strPeriod)
		if strPeriod == "" || err != nil || period < 0 {
			period = 14
		}
		df.AddRsi(period)
	}

	macd := r.URL.Query().Get("macd")
	if macd != "" {
		strPeriod1 := r.URL.Query().Get("macdPeriod1")
		strPeriod2 := r.URL.Query().Get("macdPeriod2")
		strPeriod3 := r.URL.Query().Get("macdPeriod3")
		period1, err := strconv.Atoi(strPeriod1)
		if strPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 12
		}
		period2, err := strconv.Atoi(strPeriod2)
		if strPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 26
		}
		period3, err := strconv.Atoi(strPeriod3)
		if strPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 9
		}
		df.AddMacd(period1, period2, period3)
	}

	hv := r.URL.Query().Get("hv")
	if hv != "" {
		strPeriod1 := r.URL.Query().Get("hvPeriod1")
		strPeriod2 := r.URL.Query().Get("hvPeriod2")
		strPeriod3 := r.URL.Query().Get("hvPeriod3")
		period1, err := strconv.Atoi(strPeriod1)
		if strPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 21
		}
		period2, err := strconv.Atoi(strPeriod2)
		if strPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 63
		}
		period3, err := strconv.Atoi(strPeriod3)
		if strPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 252
		}
		df.AddHv(period1)
		df.AddHv(period2)
		df.AddHv(period3)
	}

	events := r.URL.Query().Get("events")
	if events != "" {
		if config.Config.BackTest {

			df.Events = Ai.SignalEvents.CollectAfter(df.Candles[0].Time)

			/* EMAのバックテスト */
			// performance, p1, p2 := df.OptimizeEma()
			// log.Println(performance, p1, p2)
			// if performance > 0 {
			// 	df.Events = df.BackTestEma(p1, p2)
			// }

			/* ボリンジャーバンドのバックテスト */
			// performance, p1, p2 := df.OptimizeBb()
			// log.Println(performance, p1, p2)
			// if performance > 0 {
			// 	df.Events = df.BackTestBb(p1, p2)
			// }

			/* 一目均衡表のバックテスト */
			// df.Events = df.BackTestIchimoku()

			/* MACDのバックテスト */
			// performance, p1, p2, p3 := df.OptimizeMacd()
			// log.Println(performance, p1, p2, p3)
			// if performance > 0 {
			// 	df.Events = df.BackTestMacd(p1, p2, p3)
			// }

			/* RSIのバックテスト */
			// performance, period, buythread, sellthread := df.OptimizeRsi()
			// log.Println(performance, period, buythread, sellthread)
			// if performance > 0 {
			// 	df.Events = df.BackTestRsi(period, buythread, sellthread)
			// }
		} else {
			firstTime := df.Candles[0].Time // 84行めで指定した個数の頭の時間
			df.AddEvents(firstTime)
		}
	}

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
