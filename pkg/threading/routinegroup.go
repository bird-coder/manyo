/*
 * @Author: yujiajie
 * @Date: 2025-01-23 15:36:39
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 15:38:51
 * @FilePath: /Go-Base/pkg/threading/routinegroup.go
 * @Description:
 */
package threading

import "sync"

type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}

func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	})
}

func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}
