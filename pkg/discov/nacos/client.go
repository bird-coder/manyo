/*
 * @Author: yujiajie
 * @Date: 2025-03-09 15:22:17
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-15 20:01:49
 * @FilePath: /Go-Base/pkg/discov/nacos/client.go
 * @Description:
 */
package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func formatConfig(cfg *NacosConf) vo.NacosClientParam {
	clientConfig := constant.ClientConfig{
		TimeoutMs:   10000,
		NamespaceId: cfg.Client.NamespaceId,
		LogDir:      cfg.Client.LogDir,
		CacheDir:    cfg.Client.CacheDir,
		LogLevel:    cfg.Client.LogLevel,
	}
	if cfg.Client.Auth != nil {
		clientConfig.Username = cfg.Client.Auth.Username
		clientConfig.Password = cfg.Client.Auth.Password
		clientConfig.Endpoint = cfg.Client.Auth.Endpoint
		clientConfig.RegionId = cfg.Client.Auth.RegionId
		clientConfig.AccessKey = cfg.Client.Auth.AccessKey
		clientConfig.SecretKey = cfg.Client.Auth.SecretKey
	}
	serverConfigs := make([]constant.ServerConfig, len(cfg.Servers))
	for _, serverCfg := range cfg.Servers {
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(serverCfg.Ip, serverCfg.Port))
	}
	return vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	}
}

func NewNamingClient(cfg *NacosConf) (naming_client.INamingClient, error) {
	clientParams := formatConfig(cfg)
	client, err := clients.NewNamingClient(clientParams)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewConfigClient(cfg *NacosConf) (config_client.IConfigClient, error) {
	clientParams := formatConfig(cfg)
	return clients.NewConfigClient(clientParams)
}
