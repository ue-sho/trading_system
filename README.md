# 仮想通貨自動売買システム

Udemyの「現役シリコンバレーエンジニアが教えるGo入門 + 応用でビットコインのシストレFintechアプリ開発」を
受講しながら開発を行う


### config.ini

トップディレクトリにconfig.iniを作り以下の項目を設定する

```
[bitflyer]
api_key = XXXX        // APIキー
api_secret = XXXXXX   // APIの認証キー

[trading_system]
log_file = XXXXX          // logの出力先ファイル名
product_code = BTC_JPY    // 扱う仮想通貨
trade_duration = 1m       // データ取得間隔
back_test = true　        // バックテストをするかどうか
use_percent = 0.9         // 資産の90%を取引する
data_limit = 365          // 過去の取引
stop_limit_percent = 0.9  // 9割り切ったら強制的に売る
num_ranking = 3           // インディケーターの成績が良いもの上位何個使用するか

[db]
name = stockdata.sql
driver = sqlite3

[web]
port = 8080
```