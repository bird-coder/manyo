/*
 * @Author: yujiajie
 * @Date: 2025-01-24 10:03:55
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-15 21:34:13
 * @FilePath: /Go-Base/pkg/discov/nacos/subscriber.go
 * @Description:
 */
package nacos

import (
	"sync"
	"sync/atomic"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type (
	Subscriber struct {
		cfg      *NacosConf
		watchers sync.Map
	}
)

func NewSubscriber(cfg *NacosConf) *Subscriber {
	return &Subscriber{
		cfg: cfg,
	}
}

func (s *Subscriber) Watch(serviceName string) (*watcher, error) {
	item, ok := s.watchers.Load(serviceName)
	if ok {
		return item.(*watcher), nil
	}

	client, err := GetRegistry().GetNameClient(s.cfg)
	if err != nil {
		return nil, err
	}
	w, err := newWatcher(client, serviceName)
	if err != nil {
		return nil, err
	}
	s.watchers.Store(serviceName, w)
	return w, nil

}

type watcher struct {
	client           naming_client.INamingClient
	serviceName      string
	listeners        []func()
	lock             sync.RWMutex
	serviceInstances *atomic.Value
}

func newWatcher(client naming_client.INamingClient, serviceName string) (*watcher, error) {
	w := &watcher{
		client:      client,
		serviceName: serviceName,
	}

	if err := w.subscribe(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *watcher) subscribe() error {
	return w.client.Subscribe(&vo.SubscribeParam{
		ServiceName: w.serviceName,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				return
			}
			w.serviceInstances.Store(services)
			w.notify()
		},
	})
}

func (w *watcher) AddListener(listener func()) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.listeners = append(w.listeners, listener)
}

func (w *watcher) Services() []model.Instance {
	return w.serviceInstances.Load().([]model.Instance)
}

func (w *watcher) notify() {
	w.lock.RLock()
	listeners := append([]func(){}, w.listeners...)
	w.lock.RUnlock()

	for _, listener := range listeners {
		listener()
	}
}

func (w *watcher) parseInstance(instances []model.Instance) []*ServiceInstance {
	services := make([]*ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		if !instance.Healthy || !instance.Enable {
			continue
		}

		ins := &ServiceInstance{
			ServiceName: instance.ServiceName,
			ClusterName: instance.ClusterName,
			Hosts: []Host{
				{},
			},
		}

		services = append(services, ins)
	}
	return services
}
