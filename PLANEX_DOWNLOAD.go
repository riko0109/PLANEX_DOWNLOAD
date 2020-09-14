package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		log.Println("configfile is not exist!")
		bufio.NewScanner(os.Stdin).Scan()
		return
	}

	//設定を読み込む配列
	configdata, err := loadconfig(exepath + "\\config.json")
	if err != nil {
		log.Println("config decode failed:", err)
	}
	fmt.Println(configdata)

	//url := "https://svcipp.planex.co.jp/api/get_data.php?type=WS-USB01-THP&mac=24:72:60:40:03:2C&from=2020-09-04&to=2020-09-05&token=ea76a9d229497075beace88e674be9a6"

	//response, _ := http.Get(url)
	//body, _ := ioutil.ReadAll(response.Body)
	//defer response.Body.Close()

	//fmt.Println(url)
	//fmt.Println(string(body))
}

//ファイルの存在確認をする関数
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func loadconfig(path string) (*[]config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("config load failed:", err)
		return nil, err
	}
	defer file.Close()
	var configarr []config
	err = json.NewDecoder(file).Decode(&configarr)
	return &configarr, err
}

type config struct {
	DeviceName string `json:'DeviceName'`
	NicName    string `json:'NicName'`
	MacAddress string `json:'MacAddress'`
	Token      string `json:'Token'`
	SavePath   string `json:'Save_Path'`
	SaveUnit   string `json:'Save_Unit'`
}
