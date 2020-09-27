package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ue-sho/trading_system/config"
)

type SignalEvent struct {
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
}

func (s *SignalEvent) Save() bool {
	cmd := fmt.Sprintf("INSERT INTO %s (time, product_code, side, price, size) VALUES (?, ?, ?, ?, ?)", tableNameSignalEvents)
	_, err := DbConnection.Exec(cmd, s.Time.Format(time.RFC3339), s.ProductCode, s.Side, s.Price, s.Size)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") { // 同じ時間で同じ売買があったら追加しないでtrueを返す
			log.Println(err)
			return true
		}
		return false
	}
	return true
}

// BUY SELL BUY SELLを入れていく
type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

// 空を返す
func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

// 最後(最新)のシグナルイベントをloadEvents個取ってくる
func GetSignalEventsByCount(loadEvents int) *SignalEvents {
	cmd := fmt.Sprintf(`SELECT * FROM (
        SELECT time, product_code, side, price, size FROM %s WHERE product_code = ? ORDER BY time DESC LIMIT ? )
        ORDER BY time ASC;`, tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, config.Config.ProductCode, loadEvents)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return &signalEvents
}

// timeTime以降ののシグナルイベントを取ってくる
func GetSignalEventsAfterTime(timeTime time.Time) *SignalEvents {
	cmd := fmt.Sprintf(`SELECT * FROM (
                SELECT time, product_code, side, price, size FROM %s
                WHERE DATETIME(time) >= DATETIME(?)
                ORDER BY time DESC
            ) ORDER BY time ASC;`, tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, timeTime.Format(time.RFC3339))
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	return &signalEvents
}

// 空、もしくは最後がSELLかつ時間が正当なら買える
func (s *SignalEvents) CanBuy(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "SELL" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

// 空、もしくは最後がBUYかつ時間が正当なら買える
func (s *SignalEvents) CanSell(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "BUY" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

// BUYの情報
func (s *SignalEvents) Buy(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if !s.CanBuy(time) {
		return false
	}
	signalEvent := SignalEvent{
		ProductCode: ProductCode,
		Time:        time,
		Side:        "BUY",
		Price:       price,
		Size:        size,
	}
	if save { // データベースに保存する　テストするときは保存したくないので必要(バックテスト用)
		signalEvent.Save()
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Sell(productCode string, time time.Time, price, size float64, save bool) bool {

	if !s.CanSell(time) {
		return false
	}

	signalEvent := SignalEvent{
		ProductCode: productCode,
		Time:        time,
		Side:        "SELL",
		Price:       price,
		Size:        size,
	}
	if save { // データベースに保存する　テストするときは保存したくないので必要(バックテスト用)
		signalEvent.Save()
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

// SignalEvents での利益
func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SEL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding == true { // マイナスにならないように
		return beforeSell
	}
	return total
}

// JSON形式にする
func (s SignalEvents) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&struct {
		Signals []SignalEvent `json:"signals,omitempty"`
		Profit  float64       `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	})
	if err != nil {
		return nil, err
	}
	return value, err
}

// バックテスト用　データベースがないため必要
// GetSignalEventsAfterTimeと同じ
func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if time.After(signal.Time) {
			continue
		}
		return &SignalEvents{Signals: s.Signals[i:]}
	}
	return nil
}
