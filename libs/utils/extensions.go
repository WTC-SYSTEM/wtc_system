package utils

func Map[T any, K any](items []T, fn func(T) K) []K {
	var newItems []K
	for _, item := range items {
		newItems = append(newItems, fn(item))
	}
	return newItems
}

func Where[T any](items []T, fn func(T) bool) []T {
	var newItems = make([]T, 0)
	for _, item := range items {
		if fn(item) {
			newItems = append(newItems, item)
		}
	}
	return newItems
}
