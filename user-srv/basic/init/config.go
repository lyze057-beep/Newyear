package init

import (
	"5/work/Newyear/user-srv/basic/config"
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	var err error
	viper.SetConfigFile("../../dev.yaml")
	viper.ReadInConfig()
	viper.Unmarshal(&config.AppConf)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.AppConf)
}
