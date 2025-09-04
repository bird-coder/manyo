/*
 * @Author: yujiajie
 * @Date: 2025-01-24 10:03:55
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-09-04 10:37:48
 * @FilePath: /manyo/pkg/discov/nacos/subscriber.go
 * @Description:
 */
package nacos

import (
	"sync"
	"sync/atomic"
	"time"

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
	subscribeParam   *vo.SubscribeParam
	listeners        []func()
	lock             sync.RWMutex
	serviceInstances *atomic.Value
	emptyCall        atomic.Bool
	emptyTimer       *time.Timer
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
	w.subscribeParam = &vo.SubscribeParam{
		ServiceName: w.serviceName,
		SubscribeCallback: func(instances []model.Instance, err error) {
			if err != nil {
				return
			}
			if len(instances) == 0 {
				//可能是网络抖动导致获取到空列表，需要特殊处理
				if w.emptyCall.CompareAndSwap(false, true) {
					w.emptyTimer = time.AfterFunc(5*time.Second, func() {
						w.handleCallback(instances)
					})
				}
				return
			}
			w.emptyCall.Store(false)
			if w.emptyTimer != nil {
				w.emptyTimer.Stop()
				w.emptyTimer = nil
			}
			w.handleCallback(instances)
		},
	}
	return w.client.Subscribe(w.subscribeParam)
}

func (w *watcher) handleCallback(instances []model.Instance) {
	items := w.parseInstance(instances)
	w.serviceInstances.Store(items)
	w.notify()
}

func (w *watcher) AddListener(listener func()) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.listeners = append(w.listeners, listener)
}

func (w *watcher) Services() []*ServiceInstance {
	var items []*ServiceInstance
	res := w.serviceInstances.Load()
	if res != nil {
		items = res.([]*ServiceInstance)
		return items
	}
	instances, err := w.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: w.serviceName,
		HealthyOnly: true,
	})
	if err != nil {
		return nil
	}
	items = w.parseInstance(instances)
	w.serviceInstances.Store(items)
	return items
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
			ID:          instance.InstanceId,
			ServiceName: instance.ServiceName,
			ClusterName: instance.ClusterName,
			Metadata:    instance.Metadata,
			Hosts: []Host{
				{
					Ip:   instance.Ip,
					Port: instance.Port,
				},
			},
		}

		services = append(services, ins)
	}
	return services
}

func (w *watcher) Stop() error {
	err := w.client.Unsubscribe(w.subscribeParam)
	return err
}
