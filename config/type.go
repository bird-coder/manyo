/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-16 22:46:04
 * @LastEditTime: 2024-03-16 22:46:35
 * @LastEditors: yujiajie
 */
package config

type Config interface {
	LoadConfig(string) (err error)
}
