package service

import "testing"

func TestCountSIESugar(t *testing.T) {
	if err := CountSIESugar(); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestCountSIENOneBuy(t *testing.T) {
	if err := CountSIENOneBuy(); err != nil {
		t.Fatalf(err.Error())
	}
}
