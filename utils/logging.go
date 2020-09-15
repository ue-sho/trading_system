package utils

import (
	"io"
	"log"
	"os"
)

// ログファイルに書き込むように設定する
func LoggingSettings(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile err=%s", err.Error())
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)   // 複数の出力先に内容を送れる（今回の場合ログファイル）
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 日にち、時間、ファイルの何行目かがわかる
	log.SetOutput(multiLogFile)
}
