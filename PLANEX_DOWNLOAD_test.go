package main

import (
	"testing"
)

func TestExist_configが見つけられるか(t *testing.T) {
	result := Exists("C:\\Users\\Administrator\\Desktop\\PLANEX_DOWNLOAD\\config.json")
	expect := true

	if result != expect {
		t.Error("Result:", result, "Expect:", expect)
	}
	t.Log("TestExist Finished")
}
