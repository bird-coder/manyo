/*
 * @Author: yujiajie
 * @Date: 2025-01-24 11:06:47
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-15 20:02:31
 * @FilePath: /Go-Base/pkg/discov/nacos/registry.go
 * @Description:
 */
package nacos

import (
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

var (
	registry = Registry{}
)

type Registry struct {
	nameClients   sync.Map
	configClients sync.Map
}

func GetRegistry() *Registry {
	return &registry
}

func (r *Registry) GetNameClient(cfg *NacosConf) (naming_client.INamingClient, error) {
	cli, ok := r.nameClients.Load(cfg.Client.NamespaceId)
	if ok {
		return cli.(naming_client.INamingClient), nil
	}

	client, err := NewNamingClient(cfg)
	if err != nil {
		return nil, err
	}
	r.nameClients.Store(cfg.Client.NamespaceId, client)

	return client, nil
}

func (r *Registry) GetConfigClient(cfg *NacosConf) (config_client.IConfigClient, error) {
	cli, ok := r.configClients.Load(cfg.Client.NamespaceId)
	if ok {
		return cli.(config_client.IConfigClient), nil
	}

	client, err := NewConfigClient(cfg)
	if err != nil {
		return nil, err
	}
	r.configClients.Store(cfg.Client.NamespaceId, client)

	return client, nil
}
