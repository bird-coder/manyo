/*
 * @Author: yujiajie
 * @Date: 2025-01-26 14:16:38
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-15 21:09:20
 * @FilePath: /Go-Base/pkg/discov/nacos/publisher.go
 * @Description:
 */
package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type (
	Publisher struct {
		cfg *NacosConf
	}
)

func NewPublisher(cfg *NacosConf) *Publisher {
	return &Publisher{
		cfg: cfg,
	}
}

func (p *Publisher) Register(instances ...*ServiceInstance) error {
	client, err := GetRegistry().GetNameClient(p.cfg)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		for _, host := range instance.Hosts {
			_, err := client.RegisterInstance(vo.RegisterInstanceParam{
				Ip:          host.Ip,
				Port:        host.Port,
				ServiceName: instance.ServiceName,
				ClusterName: instance.ClusterName,
				GroupName:   instance.GroupName,
				Enable:      true,
				Healthy:     true,
				Ephemeral:   true,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Publisher) DeRegister(instances ...*ServiceInstance) error {
	client, err := GetRegistry().GetNameClient(p.cfg)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		for _, host := range instance.Hosts {
			_, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
				Ip:          host.Ip,
				Port:        host.Port,
				ServiceName: instance.ServiceName,
				Cluster:     instance.ClusterName,
				GroupName:   instance.GroupName,
				Ephemeral:   true,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
