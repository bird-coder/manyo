/*
 * @Author: yujiajie
 * @Date: 2025-01-23 16:15:37
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 16:29:59
 * @FilePath: /Go-Base/pkg/syncx/resourcemanager.go
 * @Description:
 */
package syncx

import (
	"io"
	"sync"
)

type ResourceManager struct {
	resources    map[string]io.Closer
	singleFlight SingleFlight
	lock         sync.RWMutex
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources:    make(map[string]io.Closer),
		singleFlight: NewSingleFlight(),
	}
}

func (manager *ResourceManager) Close() error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var errClosed error
	for _, resource := range manager.resources {
		if err := resource.Close(); err != nil {
			if errClosed == nil {
				errClosed = err
			}
		}
	}

	manager.resources = nil

	return errClosed
}

func (manager *ResourceManager) GetResource(key string, create func() (io.Closer, error)) (io.Closer, error) {
	val, err := manager.singleFlight.Do(key, func() (any, error) {
		manager.lock.RLock()
		resource, ok := manager.resources[key]
		manager.lock.RUnlock()
		if ok {
			return resource, nil
		}

		resource, err := create()
		if err != nil {
			return nil, err
		}

		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.resources[key] = resource

		return resource, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(io.Closer), nil
}

func (manager *ResourceManager) Inject(key string, resource io.Closer) {
	manager.lock.Lock()
	manager.resources[key] = resource
	manager.lock.Unlock()
}
