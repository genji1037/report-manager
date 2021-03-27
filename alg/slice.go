package alg

import (
	"bytes"
)

func SQLIn(strs []string) string {
	return "(" + StrArr2Str(strs, `'`, ",") + ")"
}

func StrArr2Str(strs []string, wrap, sep string) string {
	result := bytes.Buffer{}
	for i, str := range strs {
		if i > 0 {
			result.WriteString(sep)
		}
		result.WriteString(wrap)
		result.WriteString(str)
		result.WriteString(wrap)
	}
	return result.String()
}

func Merge() {

}
