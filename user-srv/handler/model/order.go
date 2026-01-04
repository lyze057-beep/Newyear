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
	Status       int    `gorm:"type:tinyint;default:0;comment:'商品状态:1=上架,2=下架'"`
}
type OrderMain struct {
	gorm.Model
	OrderNo     string  `gorm:"type:varchar(255);index;comment:订单编号"`
	UserID      uint    `gorm:"type:int;index;comment:关联用户id"`
	TotalAmount float64 `gorm:"type:decimal(10,2);not null;comment:订单总金额"`
	PayStatus   int     `gorm:"type:tinyint;default:0;comment:'支付状态:1=已支付,2=未支付'"`
	OrderStatus int     `gorm:"type:tinyint;default:0;comment:'订单状态:1=待发货,2=待收货,3=完成,4=取消'"`
	TradeNo     string  `gorm:"type:varchar(128);comment:支付编号"`
}
type OrderItem struct {
	gorm.Model
	OrderID      uint    `gorm:"type:int;comment:关联订单id;not null"`
	ProductID    uint    `gorm:"type:int;comment:关联商品id;not null"`
	ProductName  string  `gorm:"type:varchar(30);not null;comment:商品名称"`
	ProductPrice int     `gorm:"type:int;not null;comment:商品单价"`
	Quantity     int     `gorm:"type:int;not null;comment:购买数量"`
	Amount       float64 `gorm:"type:decimal(10,2);not null;comment:商品小计"`
}
type ShoppingCart struct {
	gorm.Model
	UserID    uint `gorm:"type:int;not null;comment:关联用户id"`
	ProductID uint `gorm:"type:int;not null;comment:关联商品id"`
	Quantity  int  `gorm:"type:int;not null;comment:购买数量"`
}
type AuthorAuth struct {
	gorm.Model
	UserID       uint   `gorm:"not null;comment:关联用户id"`
	RealName     string `gorm:"type:varchar(50);comment:实名信息"`
	AuthQualify  string `gorm:"text;comment:创作资源(比如作品链接,证明材料)"`
	Status       int    `gorm:"tinyint;default:0;comment:'审核状态 0=待审核/1=通过/2驳回'"`
	RejectReason string `gorm:"varchar(50);comment:驳回原因"`
	Auditor      string `gorm:"type:varchar(50);comment:审核人"`
}
