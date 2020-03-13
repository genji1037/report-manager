package collector

import (
	"fmt"
	"testing"
)

func TestRender(t *testing.T) {
	report := render("hello $(name), i'm $(my_name).", map[string]string{
		"name":    "James",
		"my_name": "William",
	})
	fmt.Println(report)
}
