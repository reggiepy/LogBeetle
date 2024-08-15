package arrayUtils

type Set map[string]struct{}

func NewSet(item ...string) Set {
	s := make(Set, len(item))
	for _, v := range item {
		s.Add(v)
	}
	return s
}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}
