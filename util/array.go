/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-22 22:27:43
 * @LastEditTime: 2023-09-22 23:16:01
 * @LastEditors: yuanshisan
 */
package util

type VType int

const (
	Int VType = iota
	String
	Uint32
	Uint64
	Float32
	Float64
	Bool
)

func DiffInt(src []int, dst []int) []int {
	res := src[:0]

	tmp := make(map[int]bool, len(dst))
	for _, v := range dst {
		tmp[v] = true
	}

	for _, v := range src {
		if _, ok := tmp[v]; !ok {
			res = append(res, v)
		}
	}

	return res
}

func DiffString(src []string, dst []string) []string {
	res := src[:0]

	tmp := make(map[string]bool, len(dst))
	for _, v := range dst {
		tmp[v] = true
	}

	for _, v := range src {
		if _, ok := tmp[v]; !ok {
			res = append(res, v)
		}
	}

	return res
}

func Contains(src any, val any, vType VType) (int, bool) {
	switch vType {
	case Int:
		return containsInt(src.([]int), val.(int))
	case String:
		return containsString(src.([]string), val.(string))
	case Uint32:
		return containsUint32(src.([]uint32), val.(uint32))
	case Uint64:
		return containsUint64(src.([]uint64), val.(uint64))
	case Float32:
		return containsFloat32(src.([]float32), val.(float32))
	case Float64:
		return containsFloat64(src.([]float64), val.(float64))
	}

	return -1, false
}

func containsInt(src []int, val int) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func containsString(src []string, val string) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func containsUint32(src []uint32, val uint32) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func containsUint64(src []uint64, val uint64) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func containsFloat32(src []float32, val float32) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}

func containsFloat64(src []float64, val float64) (int, bool) {
	for k, v := range src {
		if v == val {
			return k, true
		}
	}
	return -1, false
}
