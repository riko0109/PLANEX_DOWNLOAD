package main

import (
	"bufio"
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

	response, _ := http.Get(configdata[0].url())
	for i := 0; response.StatusCode != 200 && i < 10; i++ {
		body, _ := ioutil.ReadAll(response.Body)
		log.Print(string(body))
		response, _ = http.Get(configdata[0].url())
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	fmt.Println(string(body))
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

func loggingsetting(logfilepath string) {
	logfile, _ := os.OpenFile(logfilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multilogfile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
}

//URLを返却する関数
func (c *config) url() string {
	url := "https://svcipp.planex.co.jp/api/get_.php?type=" + c.DeviceName +
		"&mac=" + c.MacAddress +
		"&from=" + time.Now().AddDate(0, 0, -1).Format("2006-01-02") +
		"&to=" + time.Now().Format("2006-01-02") +
		"&token=" + c.Token
	return url
}
