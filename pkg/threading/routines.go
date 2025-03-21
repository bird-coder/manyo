/*
 * @Author: yujiajie
 * @Date: 2025-01-15 16:27:44
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 16:55:17
 * @FilePath: /manyo/pkg/threading/routines.go
 * @Description:
 */
package threading

import "github.com/bird-coder/manyo/pkg/logger"

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("runtime error: %v", err)
		}
	}()

	fn()
}
