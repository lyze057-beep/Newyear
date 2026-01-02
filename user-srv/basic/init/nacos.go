package init

import (
	"5/work/Newyear/user-srv/basic/config"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

func InitNacos() {
	nacosConf := config.AppConf.Nacos
	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConf.PublicId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            nacosConf.NaCosName,
		Password:            nacosConf.Password,
	}
	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosConf.Host,
			ContextPath: "/nacos",
			Port:        uint64(nacosConf.Port),
			Scheme:      "http",
		},
	}
	// Create config client for dynamic configuration
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(fmt.Sprintf("create config client failed, err:%v", err))
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConf.DataId,
		Group:  nacosConf.Group})
	if err != nil {
		panic(fmt.Sprintf("get config data failed, err:%v", err))
	}
	fmt.Println(string(content))
	err = yaml.Unmarshal([]byte(content), &config.AppConf)
	if err != nil {
		panic(fmt.Sprintf("unmarshal config failed, err:%v", err))
	}
}
