package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ConfigManager struct {
	cfg    *NacosConf
	client config_client.IConfigClient
}

func NewConfigManager(cfg *NacosConf) (*ConfigManager, error) {
	client, err := GetRegistry().GetConfigClient(cfg)
	if err != nil {
		return nil, err
	}
	cm := &ConfigManager{
		cfg:    cfg,
		client: client,
	}
	return cm, nil
}

func (cm *ConfigManager) GetConfig(dataId string, groupName string) (string, error) {
	content, err := cm.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupName,
	})

	return content, err
}

func (cm *ConfigManager) PutConfig(dataId string, groupName string, content string) (bool, error) {
	return cm.client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   groupName,
		Content: content,
	})
}

func (cm *ConfigManager) DelConfig(dataId string, groupName string, content string) (bool, error) {
	return cm.client.DeleteConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   groupName,
		Content: content,
	})
}

func (cm *ConfigManager) AddListener(dataId string, groupName string, fn func(string)) error {
	err := cm.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupName,
		OnChange: func(namespace, group, dataId, data string) {
			fn(data)
		},
	})
	return err
}

func (cm *ConfigManager) RemoveListener(dataId string, groupName string) error {
	err := cm.client.CancelListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupName,
	})
	return err
}

func (cm *ConfigManager) Close() {
	cm.client.CloseClient()
}
