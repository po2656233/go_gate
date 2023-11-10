package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/nothollyhigh/kiss/log"
	"go_gate/config"
	"strconv"
	"strings"
)

func getNacosClientParam() (sc []constant.ServerConfig, cc constant.ClientConfig) {
	//create ServerConfig
	sc = []constant.ServerConfig{
		*constant.NewServerConfig(config.GlobalXmlConfig.Nacos.Item.Ip, uint64(config.GlobalXmlConfig.Nacos.Item.Port), constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc = *constant.NewClientConfig(
		constant.WithNamespaceId(config.GlobalXmlConfig.Nacos.Item.Namespaceid),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(config.GlobalXmlConfig.Nacos.Logdir),
		constant.WithCacheDir(config.GlobalXmlConfig.Nacos.Cachedir), //在nacos服关闭后,使用之前缓存的内容
		constant.WithLogLevel(config.GlobalXmlConfig.Nacos.Level),
	)
	return
}

// GetHealthIP 获取有效的服务器IP
func GetHealthIP(servername string) (address string, err error) {
	// 获取参数
	sc, cc := getNacosClientParam()

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Error("NewNamingClient err:%v ", err.Error())
		return "", err
	}

	clusters := strings.Split(config.GlobalXmlConfig.Nacos.Item.Clusters, ",")

	//// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	//instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
	//	ServiceName: servername,
	//	GroupName:   config.GlobalXmlConfig.Nacos.Item.Groupname,             // 默认值DEFAULT_GROUP
	//	Clusters:    clusters, // 默认值DEFAULT
	//})
	//if err != nil{
	//	log.Error("SelectAllInstances err:%v ", err.Error())
	//	return "", err
	//}
	//fmt.Println(instances)

	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: servername,
		GroupName:   config.GlobalXmlConfig.Nacos.Item.Groupname, // 默认值DEFAULT_GROUP
		Clusters:    clusters,                                    // 默认值DEFAULT
	})
	if err != nil {
		log.Error("SelectOneHealthyInstance err:%v ", err.Error())
		return "", err
	}
	addr := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
	return addr, err
}

// Subscribe 订阅处理服务列表
func Subscribe(servername string, f func(services []model.SubscribeService, err error)) error {
	// 获取参数
	sc, cc := getNacosClientParam()
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Error("NewNamingClient err:%v ", err.Error())
		return err
	}

	clusters := strings.Split(config.GlobalXmlConfig.Nacos.Item.Clusters, ",")
	err = namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName:       servername,
		GroupName:         config.GlobalXmlConfig.Nacos.Item.Groupname, // 默认值DEFAULT_GROUP
		Clusters:          clusters,                                    // 默认值DEFAULT
		SubscribeCallback: f,
	})
	return err
}
