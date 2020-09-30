package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"planexdownload/config"
	"sync"
)

const (
	configfilename string = "config.yaml" //設定ファイルのファイル名
	logfilename    string = "Log.log"     //ログファイルのファイル名
)

func main() {
	//実行ファイルのフルパスを取得
	exepath, err := os.Executable()
	if err != nil {
		log.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	//実行ファイルがあるディレクトリを取得
	exepath = filepath.Dir(exepath)

	//設定ファイル存在確認
	if !config.Exists(filepath.Join(exepath, configfilename)) {
		log.Println(filepath.Join(exepath+configfilename) + "が存在しません")
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	//ログ書き出し設定
	loggingsetting(filepath.Join(exepath, logfilename))

	//設定を読み込む配列
	configdata, err := config.Loadconfig(filepath.Join(exepath, configfilename))
	if err != nil {
		log.Fatalln(err)
	}
	//goroutine待ち合わせ用構造体
	var wg sync.WaitGroup
	//並列処理でAPI叩いてCSV書き出しまで
	for i := 0; i < len(configdata); i++ {
		//処理開始時にwgのカウント加算
		wg.Add(1)
		go func(i2 int) {
			err := configdata[i2].Getstringarrayfromapi()
			if err != nil {
				log.Panicln(err)
			}
			err = configdata[i2].Createcsv()
			if err != nil {
				log.Panicln(err)
			}
			//処理終了時にwgの処理減算
			defer wg.Done()
		}(i)
	}
	//wgのカウントが0になるまで待機
	wg.Wait()
	fmt.Println("コンソールを終了するには何かキーを押してください…")
	bufio.NewScanner(os.Stdin).Scan()
}

//ログを吐き出すもろもろの設定をする関数
func loggingsetting(logfilepath string) {
	//ログファイルを開く
	//os.O_RDWR=読書FLAG
	//os.O_CREATE=存在しなかったら生成
	//os.O_APPEND=存在したら追記
	//0666=パーミッション
	logfile, _ := os.OpenFile(logfilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//ログファイルを標準出力とテキストファイルに出力する設定
	multilogfile := io.MultiWriter(os.Stdout, logfile)
	//出力する項目を設定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//設定を実行
	log.SetOutput(multilogfile)
}
