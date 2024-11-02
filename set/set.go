package set

// Set represents a generic set
type Set[T comparable] map[T]struct{}

func NewWithCap[T comparable](cap int) Set[T] {
	return make(Set[T], cap)
}

func (set Set[T]) Add(element T) {
	set[element] = struct{}{}
}

func (set Set[T]) AddSlice(elements []T) {
	for _, element := range elements {
		set[element] = struct{}{}
	}
}

func (set Set[T]) AddSet(other Set[T]) {
	for element := range other {
		set[element] = struct{}{}
	}
}

func (set Set[T]) Remove(element T) {
	delete(set, element)
}

func (set Set[T]) RemoveSlice(elements []T) {
	for _, element := range elements {
		delete(set, element)
	}
}

func (set Set[T]) RemoveSet(other Set[T]) {
	for element := range other {
		delete(set, element)
	}
}

func (set Set[T]) Contains(element T) bool {
	_, exists := set[element]
	return exists
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) IsEmpty() bool {
	return len(set) == 0
}

func (set Set[T]) Clear() {
	for element := range set {
		delete(set, element)
	}
}

func (set Set[T]) Elements() []T {
	elements := make([]T, 0, len(set))
	for element := range set {
		elements = append(elements, element)
	}
	return elements
}
