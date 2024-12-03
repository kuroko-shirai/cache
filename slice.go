package cache

type slice[K comparable] []K

func (s slice[K]) find(key K) (int, bool) {
	for i, k := range s {
		if k == key {
			return i, true
		}
	}

	return 0, false
}

func (s *slice[K]) delete(key K) bool {
	i, ok := s.find(key)
	if !ok {
		return false
	}

	(*s)[i] = (*s)[len((*s))-1]
	(*s)[len((*s))-1] = *new(K)
	(*s) = (*s)[:len((*s))-1]

	return true
}

func (s *slice[K]) append(key K) {
	*s = append(*s, key)
}
