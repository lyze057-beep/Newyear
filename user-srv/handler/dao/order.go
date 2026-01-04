package dao

import (
	"5/work/Newyear/user-srv/basic/config"
	"5/work/Newyear/user-srv/handler/model"
	"errors"

	"gorm.io/gorm"
)

type Register struct{}
type Login struct{}
type UpdatePassword struct{}
type ListDao struct{}
type Product struct{}
type OrderDao struct{}
type CartDao struct{}
type AuthorAuthDao struct{}

// 用户添加
func (r *Register) Create(user *model.User) error {
	return config.DB.Create(&user).Error
}

// 用户查询
func (r *Register) GetUser(username string) (*model.User, error) {
	var u model.User
	if err := config.DB.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// 修改密码
func (u *UpdatePassword) UpdateDao(username, newPassword string) error {
	result := config.DB.Model(&model.User{}).Where("username=?", username).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// 查询用户所有
func (l *ListDao) ListUser(db *gorm.DB) ([]model.User, error) {
	var users []model.User
	config.DB.Find(&users)
	return users, nil
}

// 添加商品
func (p *Product) CreateProduct(pro *model.Product) error {
	return config.DB.Create(pro).Error
}

// 通过id查询商品
func (p *Product) Get(id uint) (*model.Product, error) {
	var pro model.Product
	if err := config.DB.Where("id=?", id).First(&pro).Error; err != nil {
		return nil, err
	}
	return &pro, nil
}

// 获取商品列表倒序
func (p *Product) List() ([]model.Product, error) {
	var list []model.Product
	if err := config.DB.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// 添加订单
func (o *OrderDao) OrderAddDao(der *model.OrderMain) error {
	return config.DB.Create(&der).Error
}

// MarkPaid 标记订单支付成功并写入支付宝交易号
// pay_status=1 表示已支付
func (o *OrderDao) MarkPaid(orderNo string, tradeNo string) error {
	return config.DB.Model(&model.OrderMain{}).Where("order_no=?", orderNo).
		Updates(map[string]interface{}{
			"pay_status": 1,
			"trade_no":   tradeNo,
		}).Error
}

// 通过支付状态查询
func (o *OrderDao) GetOrder(PayStatus int) (*model.OrderMain, error) {
	var order model.OrderMain
	if err := config.DB.Where("pay_status=?", PayStatus).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// 购物车添加条目，存在则累加数量
func (d *CartDao) Add(userID uint, productID uint, qty int) error {
	var item model.ShoppingCart
	// 尝试查询是否存在
	err := config.DB.Where("user_id=? and product_id=?", userID, productID).First(&item).Error

	if err == nil {
		// 存在则更新数量（这里选择覆盖，如果需要累加可以用 +=）
		item.Quantity = qty
		return config.DB.Save(&item).Error
	}

	// 不存在则创建新条目
	// 注意：这里必须使用传入参数 userID 和 productID，不能使用 item.UserID（因为 item 是空的）
	newItem := model.ShoppingCart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  qty,
	}
	return config.DB.Create(&newItem).Error
}

// list获取用户购物车列表
func (d *CartDao) List(userID uint) ([]model.ShoppingCart, error) {
	var items []model.ShoppingCart
	if err := config.DB.Where("user_id=?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Remove从购物车中移除指定的商品
func (d *CartDao) Remove(userID uint, productID uint) error {
	return config.DB.Where("user_id=? and product_id=?", userID, productID).Delete(&model.ShoppingCart{}).Error
}

// Clear清空购物车
func (d *CartDao) Clear(userID uint) error {
	return config.DB.Where("user_id=?", userID).Delete(&model.ShoppingCart{}).Error
}

// 提交认证申请
func (a *AuthorAuthDao) Create(auth *model.AuthorAuth) error {
	return config.DB.Create(&auth).Error
}

// 查询认证状态
func (a *AuthorAuthDao) Get(userID uint) (*model.AuthorAuth, error) {
	var auth model.AuthorAuth
	if err := config.DB.Where("user_id=?", userID).First(&auth).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

// 审核状态
func (a *AuthorAuthDao) Update(authID uint, status int, rejectReason, auditor string) error {
	return config.DB.Model(&model.AuthorAuth{}).Where("id=?", authID).
		Updates(map[string]interface{}{
			"status":        status,
			"reject_reason": rejectReason,
			"auditor":       auditor,
		}).Error
}
