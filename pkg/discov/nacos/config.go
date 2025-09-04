/*
 * @Author: yujiajie
 * @Date: 2025-01-24 09:06:26
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-09-03 18:14:03
 * @FilePath: /manyo/pkg/discov/nacos/config.go
 * @Description:
 */
package nacos

type NacosConf struct {
	Client  NacosClientConf   `mapstructure:"client"`
	Servers []NacosServerConf `mapstructure:"servers"`
}

type NacosClientConf struct {
	NamespaceId string     `mapstructure:"namespace"`
	LogDir      string     `mapstructure:"logdir"`
	CacheDir    string     `mapstructure:"cachedir"`
	LogLevel    string     `mapstructure:"loglevel"`
	Auth        *NaocsAuth `mapstructure:"auth"`
}

type NaocsAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"pwd_password"`

	//以下是使用阿里云nacos时配置
	Endpoint  string `mapstructure:"endpoint"` //阿里云服务端点，配置后可以不需要 NacosServerConf
	RegionId  string `mapstructure:"regionId"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"pwd_secretKey"`
}

type NacosServerConf struct {
	Ip   string `mapstructure:"ip"`
	Port uint64 `mapstructure:"port"`
}

type Host struct {
	Ip   string `mapstructure:"ip"`
	Port uint64 `mapstructure:"port"`
}

type ServiceInstance struct {
	ID          string            //服务编号
	ServiceName string            //服务名
	ClusterName string            //机器集群名
	GroupName   string            //业务分组名
	Metadata    map[string]string //服务原数据
	Hosts       []Host
}
