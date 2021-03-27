package alg

type StrMap struct {
	m map[string]struct{}
}

func NewStrMapFromSlice(strs []string) *StrMap {
	m := make(map[string]struct{})
	for _, str := range strs {
		m[str] = struct{}{}
	}
	return &StrMap{m: m}
}

func (s *StrMap) Contain(key string) bool {
	_, ok := s.m[key]
	return ok
}
