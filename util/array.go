/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-22 22:27:43
 * @LastEditTime: 2024-03-11 21:31:02
 * @LastEditors: yujiajie
 */
package util

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
