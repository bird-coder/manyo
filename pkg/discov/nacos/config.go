/*
 * @Author: yujiajie
 * @Date: 2025-01-24 09:06:26
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-09 16:42:37
 * @FilePath: /Go-Base/pkg/discov/nacos/config.go
 * @Description:
 */
package nacos

type NacosConf struct {
	Client  NacosClientConf
	Servers []NacosServerConf
}

type NacosClientConf struct {
	NamespaceId string
	LogDir      string
	CacheDir    string
	LogLevel    string
	Auth        *NaocsAuth
}

type NaocsAuth struct {
	Username string
	Password string

	//以下是使用阿里云nacos时配置
	Endpoint  string //阿里云服务端点，配置后可以不需要 NacosServerConf
	RegionId  string
	AccessKey string
	SecretKey string
}

type NacosServerConf struct {
	Host
}

type Host struct {
	Ip   string
	Port uint64
}

type ServiceInstance struct {
	ServiceName string //服务名
	ClusterName string //机器集群名
	GroupName   string //业务分组名
	Hosts       []Host
}
