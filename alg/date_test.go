package alg

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println(time.Now().In(time.UTC).Format("2006-01-02 15:04:05"))
}
