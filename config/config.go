package config

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	url string = "https://svcipp.planex.co.jp/api/get_data.php?type=" //APIのURL
)

//Config is JSONから読み込んだ設定の構造体
type Config struct {
	DeviceName string     `json:"DeviceName"` //デバイス名 WS-USB01-THPとか
	NicName    string     `json:"NicName"`    //名前　サーバールームとか
	MacAddress string     `json:"MacAddress"` //MACアドレス
	Token      string     `json:"Token"`      //トークン
	SavePath   string     `json:"SavePath"`   //CSVを保存するディレクトリ
	From       int        `json:"From"`       //何日前から 実行日から1日前なら-1
	To         int        `json:"To"`         //何日目迄取得するか　実行日までなら0
	GetData    [][]string //APIから返却されてきたJSON用のジャグ配列
}

//Configyaml is yamlから読み込んだ設定の構造体
type Configyaml struct {
	DeviceName string     `yaml:"DeviceName"` //デバイス名 WS-USB01-THPとか
	NicName    string     `yaml:"NicName"`    //名前　サーバールームとか
	MacAddress string     `yaml:"MacAddress"` //MACアドレス
	Token      string     `yaml:"Token"`      //トークン
	SavePath   string     `yaml:"SavePath"`   //CSVを保存するディレクトリ
	From       int        `yaml:"From"`       //何日前から 実行日から1日前なら-1
	To         int        `yaml:"To"`         //何日目迄取得するか　実行日までなら0
	GetData    [][]string //APIから返却されてきたJSON用のジャグ配列
}

//URL is 組み立てたURLを返却する関数
func (c *Configyaml) URL() string {
	url := url + c.DeviceName +
		"&mac=" + c.MacAddress +
		"&from=" + time.Now().AddDate(0, 0, c.From).Format("2006-01-02") +
		"&to=" + time.Now().AddDate(0, 0, c.To+1).Format("2006-01-02") +
		"&token=" + c.Token
	return url
}

//Loadconfig is configを読み込み構造体に変換する関数
func Loadconfig(path string) ([]Configyaml, error) {
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

	var configs []Configyaml
	//jsonを構造体のスライスに変換
	err = yaml.UnmarshalStrict(b, &configs)
	if err != nil {
		return nil, err
	}
	return configs, err
}

//Getstringarrayfromapi is
//API叩いてデータを取得して文字列型に変換して構造体のGetDataに格納する関数
func (c *Configyaml) Getstringarrayfromapi() error {
	fmt.Println(c.NicName + ":ダウンロード開始!")
	//APIにGETリクエスト
	response, _ := http.Get(c.URL())
	//レスポンスのステータスが200になるまで20回繰り返す
	for i := 0; i < 20 && response.StatusCode != 200; i++ {
		log.Println(c.NicName + ":" + response.Status + "リトライ")
		response, _ = http.Get(c.URL())
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

//Createcsv is csvを生成してデータを書き込む関数
func (c *Configyaml) Createcsv() error {
	if !Exists(filepath.Join(c.SavePath, time.Now().AddDate(0, 0, c.From).Format("200601"))) {
		err := os.Mkdir(filepath.Join(c.SavePath, time.Now().AddDate(0, 0, c.From).Format("200601")), 0777)
		if err != nil {
			return err
		}
	}
	c.SavePath = filepath.Join(c.SavePath, time.Now().AddDate(0, 0, c.From).Format("200601"))
	//ファイルネームは名前+FROM+TO
	filename := filepath.Join(c.SavePath, c.NicName+"_"+
		time.Now().AddDate(0, 0, c.From).Format("20060102")+"_"+
		time.Now().AddDate(0, 0, c.To).Format("20060102")+".csv")

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

//Exists is ファイルがあるかどうかを調べる関数
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
