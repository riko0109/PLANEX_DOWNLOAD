package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	configfilename string = "config.json"                                        //設定ファイルのファイル名
	logfilename    string = "Log.log"                                            //ログファイルのファイル名
	url            string = "https://svcipp.planex.co.jp/api/get_data.php?type=" //APIのURL
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
	if !exists(filepath.Join(exepath, configfilename)) {
		log.Println(filepath.Join(exepath+configfilename) + "が存在しません")
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	//ログ書き出し設定
	loggingsetting(filepath.Join(exepath, logfilename))

	//設定を読み込む配列
	configdata, err := loadconfig(filepath.Join(exepath, configfilename))
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
			err := configdata[i2].getstringarrayfromapi()
			if err != nil {
				log.Panicln(err)
			}
			err = configdata[i2].createcsv()
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

//ファイルがあるかどうかを調べる関数
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//コンフィグJSONを読み込み構造体に変換する関数
func loadconfig(path string) ([]config, error) {
	//ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//ファイルを読み取る
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configs []config
	//jsonを構造体のスライスに変換
	err = json.Unmarshal(b, &configs)
	if err != nil {
		return nil, err
	}
	return configs, err
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

//JSONから読み込んだ設定の構造体
type config struct {
	DeviceName string     `json:"DeviceName"` //デバイス名 WS-USB01-THPとか
	NicName    string     `json:"NicName"`    //名前　サーバールームとか
	MacAddress string     `json:"MacAddress"` //MACアドレス
	Token      string     `json:"Token"`      //トークン
	SavePath   string     `json:"SavePath"`   //CSVを保存するディレクトリ
	From       int        `json:"From"`       //何日前から 実行日から1日前なら-1
	To         int        `json:"To"`         //何日目迄取得するか　実行日までなら0
	GetData    [][]string //APIから返却されてきたJSON用のジャグ配列
}

//URLを返却する関数
func (c *config) url() string {
	url := url + c.DeviceName +
		"&mac=" + c.MacAddress +
		"&from=" + time.Now().AddDate(0, 0, c.From).Format("2006-01-02") +
		"&to=" + time.Now().AddDate(0, 0, c.To+1).Format("2006-01-02") +
		"&token=" + c.Token
	return url
}

//API叩いてデータを取得して文字列型に変換して構造体のGetDataに格納する関数
func (c *config) getstringarrayfromapi() error {
	fmt.Println(c.NicName + ":ダウンロード開始!")
	//APIにGETリクエスト
	response, _ := http.Get(c.url())
	//レスポンスのステータスが200になるまで20回繰り返す
	for i := 0; i < 20 && response.StatusCode != 200; i++ {
		log.Println(c.NicName + ":" + response.Status + "リトライ")
		response, _ = http.Get(c.url())
	}
	//20回のリクエストでレスポンスが無かったらタイムアウト
	if response.StatusCode != 200 {
		return errors.New("APIからのレスポンスがタイムアウトしました")
	}
	//レスポンスのボディを読み取り
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	//GetDataにジャグ配列として格納
	if err := json.Unmarshal(body, &c.GetData); err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.NicName + ":ダウンロード完了!")
	return nil
}

//csvを生成してデータを書き込む関数
func (c *config) createcsv() error {
	//ファイルネームは名前+FROM+TO
	filename := c.SavePath + "\\" + c.NicName + "_" +
		time.Now().AddDate(0, 0, c.From).Format("20060102") + "_" +
		time.Now().AddDate(0, 0, c.To).Format("20060102") + ".csv"

	//ファイルを開く
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	//関数終了時にクローズ
	defer file.Close()

	//CSV書き込み
	writer := csv.NewWriter(file)
	writer.UseCRLF = true
	//バッファにためてから
	writer.WriteAll(c.GetData)
	writer.Flush()
	//バッファ内を書き込み
	fmt.Println(c.NicName + ": csv書き込み完了!")
	return nil
}
