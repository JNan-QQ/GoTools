/*
	基本类型切片相关操作
*/

package data

import "slices"

// IsEmpty 判断切片是否为空
func IsEmpty(array []string) bool {
	if len(array) == 0 {
		return true
	}
	for _, s := range array {
		if s != "" {
			return false
		}
	}
	return true
}

// RepeatIndex 提取切片内重复值的索引
func RepeatIndex[E comparable](array []E) map[E][]int {

	result := map[E][]int{}

	for i, e := range array {
		if _, ok := result[e]; ok {
			result[e] = append(result[e], i)
		} else {
			result[e] = []int{i}
		}
	}

	for key, value := range result {
		if len(value) == 1 {
			delete(result, key)
		}
	}

	return result
}

// Filter 切片过滤
func Filter[T comparable](slice []T, condition func(T) bool) ([]T, []int) {
	var filtered []T
	var indexes []int
	for index, item := range slice {
		if condition(item) {
			filtered = append(filtered, item)
			indexes = append(indexes, index)
		}
	}
	return filtered, indexes
}

func SliceToAny[T comparable](slice []T) []any {
	var anySlice []any
	for _, item := range slice {
		anySlice = append(anySlice, item)
	}
	return anySlice
}

func CreateSlice[E any](v E, l int) []E {
	var values []E
	for i := 0; i < l; i++ {
		values = append(values, v)
	}
	return values
}

// Range 生成序列切片
func Range(start, stop, step int) []int {
	var x []int
	for i := start; i < stop; i += step {
		x = append(x, i)
	}
	return x
}

func Overlap[T comparable](s1, s2 []T) []T {
	var overlap []T
	for _, t := range s1 {
		if slices.Contains(s2, t) {
			overlap = append(overlap, t)
		}
	}
	return overlap
}
