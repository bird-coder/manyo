/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-06 22:49:35
 * @LastEditTime: 2024-03-06 22:53:56
 * @LastEditors: yujiajie
 */
package constant

type Code int

const (
	SUCCESS Code = 100 + iota
	ERROR
)

type ErrorMsg string
