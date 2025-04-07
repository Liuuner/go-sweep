package common

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}
func (s *Set[T]) Add(value T) {
	(*s)[value] = struct{}{}
}
func (s *Set[T]) remove(value T) {
	delete(*s, value)
}
func (s Set[T]) intersection(otherSet Set[T]) []T {
	result := make([]T, 0)
	for v := range s {
		if otherSet.Has(v) {
			result = append(result, v)
		}
	}
	return result
}
