package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	if !Exists(exepath + "PLANEX_DOWNLOAD.conf") {
		log.Println("conf is not exist!")
		bufio.NewScanner(os.Stdin).Scan()
		return
	}

	//設定を読み込む連想配列
	configdata := make(map[string]string, 100)

	url := "https://svcipp.planex.co.jp/api/get_data.php?type=WS-USB01-THP&mac=24:72:60:40:03:2C&from=2020-09-04&to=2020-09-05&token=ea76a9d229497075beace88e674be9a6"

	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	fmt.Println(url)
	fmt.Println(string(body))
}

//ファイルの存在確認をする関数
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
