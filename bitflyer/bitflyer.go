package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

// BiyFlyerでAPI認証を行う
// BitFlyerの認証の参考例をGoで書き換えた (https://lightning.bitflyer.com/docs)
func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)   // timestampを10進数にする
	message := timestamp + method + endpoint + string(body) // ACCESS-TIMESTAMP, HTTP メソッド, リクエストのパス, リクエストボディ

	mac := hmac.New(sha256.New, []byte(api.secret)) // Keyed-Hash Message Authentication Code (HMAC)
	mac.Write([]byte(message))                      // api.secretを鍵として、messageを暗号化
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

// Httpリクエストを行う
func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}
	endpoint := baseURL.ResolveReference(apiURL).String()
	log.Printf("action=doRequest endpoint=%s", endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data)) // POSTの場合は送るデータを第３引数にいれる
	if err != nil {
		return
	}
	q := req.URL.Query()
	for key, value := range query { // クエリを付け加える
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode() // map[a:[1] b:[2] c[3&$] を　a=1&b=2&c=3%26%25　というようにする

	for key, value := range api.header(method, req.URL.RequestURI(), data) { // ヘッダ情報を付け加える（認証）
		req.Header.Add(key, value)
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

// BitFlyerにある資産データを取得する
func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	log.Printf("url=%s resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	return balance, nil
}
