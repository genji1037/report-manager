package proxy

import (
	"fmt"
	"testing"
)

func TestGetPlatformSnapshot(t *testing.T) {

	got, err := GetPlatformSnapshot("2021-05-26")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(got)

}
