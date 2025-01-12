package utils

import "github.com/JuanGQCadavid/r-tree/app/core/domain"

func DeleteElements[T any](slice []*domain.Location[T], index ...int) []*domain.Location[T] {
	result := make([]*domain.Location[T], len(slice)-len(index))
	counter := 0

	for i := range slice {
		toAdd := true

		for _, v := range index {
			if v == i {
				toAdd = false
				break
			}
		}

		if toAdd {
			result[counter] = slice[i]
			counter++
		}
	}
	return result
}
