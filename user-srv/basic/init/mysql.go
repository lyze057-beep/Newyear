package init

import (
	"5/work/Newyear/user-srv/basic/config"
	"5/work/Newyear/user-srv/handler/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	var err error
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	mysqlConf := config.AppConf.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/newyear?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port)
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql init success")
	config.DB.AutoMigrate(&model.User{}, &model.Product{}, &model.OrderMain{},
		&model.OrderItem{}, &model.ShoppingCart{})
}
