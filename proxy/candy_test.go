package proxy

import (
	"fmt"
	"report-manager/config"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.LoadConfig("../config.yaml")
	if err != nil {
		panic(m)
	}
	m.Run()
}

func TestLatestCirculateAmount(t *testing.T) {
	ca, err := LatestCirculateAmount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ca.String())
}
