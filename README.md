# 仮想通貨自動売買システム

Udemyの「現役シリコンバレーエンジニアが教えるGo入門 + 応用でビットコインのシストレFintechアプリ開発」を
受講しながら開発を行う


### config.ini

topディレクトリにconfig.iniを作り以下の項目を設定する

[bitflyer]
api_key = XXXX // APIキー
api_secret = XXXXXX // APIの認証キー

[trading_system]
log_file = XXXXX // logの出力先ファイル名
product_code = BTC_JPY