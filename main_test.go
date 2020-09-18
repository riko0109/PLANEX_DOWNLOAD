package main

import (
	"testing"
)

func TestExistssuccess(t *testing.T) { //exist正常形
	result := exists("C:")
	expect := true

	if result != expect {
		t.Error("Result:", result, "Expect:", expect)
	}
	t.Log("Testexistsuccess:OK")
}

func TestExistsfail(t *testing.T) { //exist異常形
	result := exists("a:")
	expect := false

	if result != expect {
		t.Error("Result:", result, "Expect:", expect)
	}
	t.Log("Testexistsuccess:OK")
}
