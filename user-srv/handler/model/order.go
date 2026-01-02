package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(30);not null;comment:用户名"`
	Password string `gorm:"type:varchar(64);not null;comment:密码"`
	Phone    string `gorm:"type:char(11);not null;comment:手机号"`
}
type Product struct {
	gorm.Model
	ProductName  string `gorm:"type:varchar(30);not null;comment:商品名称"`
	ProductPrice int    `gorm:"type:int;not null;comment:商品单价"`
	ProductNum   int    `gorm:"type:int;not null;comment:商品库存"`
}
type Order struct {
	gorm.Model
	UserID uint `gorm:"type:"`
}
