/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-11 22:13:43
 * @LastEditTime: 2024-03-11 22:46:47
 * @LastEditors: yujiajie
 */
package response

type Response interface {
	SetCode(int)
	SetMsg(string)
	SetData(interface{})
	Clone() Response
}
