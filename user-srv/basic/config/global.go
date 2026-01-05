package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	AppConf AppConfig
	DB      *gorm.DB
	Rdb     *redis.Client
)
