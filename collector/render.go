package collector

import "strings"

func render(ori string, kv map[string]string) string {
	out := ori
	for k, v := range kv {
		variable := "$(" + k + ")"
		out = strings.ReplaceAll(out, variable, v)
	}
	return out
}

type StringRender struct {
	label string
	str   string
}

func (s *StringRender) Collect() error {
	return nil
}

func (s *StringRender) Render(ori string) string {
	return render(ori, map[string]string{
		s.label: s.str,
	})
}

func NewStringRender(label, str string) *StringRender {
	return &StringRender{str: str, label: label}
}
