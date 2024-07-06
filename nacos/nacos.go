package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"thundersoft.com/brain/DigitalVisitor/conf"
)

func GetKnowlede() (string, uint64) {
	// 创建 Nacos å®¢户端配置
	sc := []constant.ServerConfig{
		{
			IpAddr: conf.ConfigInfo.Nacos.Host,
			Port:   uint64(conf.ConfigInfo.Nacos.Port),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         "rubik-brain-dev",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
	}

	// 创建 Nacos 客户端
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		fmt.Println("Error creating Nacos client:", err)
		return "", 0
	}

	// 获取服务信息
	serviceName := "chat-knowledge-service"
	groupName := "DEFAULT_GROUP"
	instances, err := client.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	if err != nil {
		fmt.Println("Error getting service instances:", err)
		return "", 0
	}

	// 处理服务信息
	for _, instance := range instances {
		fmt.Println("Instance:", instance.Ip, instance.Port)
		return instance.Ip, instance.Port
	}
	return "", 0
}
