/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-22 22:27:43
 * @LastEditTime: 2024-12-27 16:10:46
 * @LastEditors: yujiajie
 */
package array

// 比较两个切片，获取第一个切片相对于第一个切片的差集
func DiffSlice[S ~[]E, E comparable](src S, dst S) []E {
	res := src[:0]

	tmp := make(map[E]bool, len(dst))
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

// 比较两个切片，取交集
func IntersectSlice[S ~[]E, E comparable](src S, dst S) []E {
	res := src[:0]

	tmp := make(map[E]bool, len(dst))
	for _, v := range dst {
		tmp[v] = true
	}

	for _, v := range src {
		if _, ok := tmp[v]; ok {
			res = append(res, v)
		}
	}

	return res
}
