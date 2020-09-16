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
	"time"
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
	if !Exists(exepath + "\\config.json") {
		log.Println("configfile is not exist!" + exepath + "\\config.json")
		bufio.NewScanner(os.Stdin).Scan()
		os.Exit(1)
	}

	loggingsetting(exepath + "\\HTTP_Failed.log")

	//設定を読み込む配列
	configdata, err := loadconfig(exepath + "\\config.json")
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < len(configdata); i++ {
		err := configdata[i].getstringarrayfromapi()
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println(configdata[0].url())
	//fmt.Println(configdata[0].GetData)
	configdata[0].createcsv()
	bufio.NewScanner(os.Stdin).Scan()
}

//ファイルの存在確認をする関数
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//コンフィグJSONを読み込み構造体に変換する関数
func loadconfig(path string) ([]config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("config open failed! : ")
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("config load failed : ")
	}
	defer file.Close()

	var configs []config

	err = json.Unmarshal(b, &configs)
	if err != nil {
		return nil, errors.New("json→struct convert failed : ")
	}
	return configs, err
}

//ログを吐き出すもろもろの設定をする関数
func loggingsetting(logfilepath string) {
	//ログファイルを開く
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
	DeviceName string `json:"DeviceName"`
	NicName    string `json:"NicName"`
	MacAddress string `json:"MacAddress"`
	Token      string `json:"Token"`
	SavePath   string `json:"SavePath"`
	SaveUnit   string `json:"SaveUnit"`
	Url        string
	GetData    [][]string
}

type jsontype [][]string

//URLを返却する関数
func (c *config) url() string {
	url := "https://svcipp.planex.co.jp/api/get_data.php?type=" + c.DeviceName +
		"&mac=" + c.MacAddress +
		"&from=" + time.Now().AddDate(0, 0, -1).Format("2006-01-02") +
		"&to=" + time.Now().Format("2006-01-02") +
		"&token=" + c.Token
	return url
}

//API叩いてデータを取得して文字列型に変換して構造体のGetDataに格納する関数
func (c *config) getstringarrayfromapi() error {
	response, _ := http.Get(c.url())

	for i := 0; i < 10 && response.StatusCode != 200; i++ {
		body, _ := ioutil.ReadAll(response.Body)
		log.Println(string(body))
		response, _ = http.Get(c.url())
	}
	if response.StatusCode != 200 {
		return errors.New("apireturn is timeout!")
	}
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err := json.Unmarshal(body, &c.GetData); err != nil {
		log.Fatal(err)
	}
	return nil
}

//csvを生成してデータを書き込む関数
func (c *config) createcsv() error {
	file, err := os.OpenFile(c.SavePath+"\\"+c.NicName+time.Now().AddDate(0, 0, -1).Format("20060102")+".csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.New("File Create Failed!")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.UseCRLF = true
	writer.WriteAll(c.GetData)
	writer.Flush()
	return nil
}
