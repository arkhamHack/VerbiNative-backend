package websockets

func Except[T any](arr []T, predicate func(item T) bool) []T {
	var pos = -1
	for i, item := range arr {
		if predicate(item) {
			pos = i
			break
		}
	}
	if pos != -1 {
		result := make([]T, 0, len(arr)-1)
		result = append(result, arr[:pos]...)
		result = append(result, arr[pos+1:]...)
		return result

	}
	return arr
}

func ForEach[T any](arr []T, method func(item T)) {
	for _, item := range arr {
		method(item)
	}
}
