package config

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	AppConf AppConfig
	DB      *gorm.DB
	Rdb     *redis.Client
)
