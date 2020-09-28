package config

import (
	"testing"
	"time"
)

func TestUrlsuccess(t *testing.T) { //url正常形
	testconfig := Configyaml{DeviceName: "DeviceName",
		NicName:    "",
		MacAddress: "MacAddress",
		Token:      "abcdefghijklmn",
		SavePath:   "",
		From:       0,
		To:         0,
		GetData:    nil,
	}
	result := "https://svcipp.planex.co.jp/api/get_data.php?type=DeviceName&mac=MacAddress&from=" +
		time.Now().AddDate(0, 0, testconfig.From).Format("2006-01-02") + "&to=" +
		time.Now().AddDate(0, 0, testconfig.To+1).Format("2006-01-02") + "&token=abcdefghijklmn"
	expect := testconfig.URL()

	if result != expect {
		t.Error("result:", result, "expect:", expect)
	}
	t.Log("TestUrlsuccess：OK")
}