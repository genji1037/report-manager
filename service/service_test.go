package service

import (
	"report-manager/config"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.LoadConfig("../config.yaml")
	if err != nil {
		panic(err)
	}
	m.Run()
}
