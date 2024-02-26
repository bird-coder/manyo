/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-12 22:16:30
 * @LastEditTime: 2023-10-14 20:20:45
 * @LastEditors: yuanshisan
 */
package constant

type Env string

const (
	Dev  Env = "development"
	Prod Env = "production"
)

func (e Env) String() string {
	return string(e)
}
