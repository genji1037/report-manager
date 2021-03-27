package alg

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteTemp(t *testing.T) {

	rd := bytes.NewBufferString("hello")
	f, err := WriteTemp(rd, "fname.txt")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(f.Name())
	f.Close()

}
