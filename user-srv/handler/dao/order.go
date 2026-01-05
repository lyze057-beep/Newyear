package dao

import (
	"5/work/Newyear/user-srv/basic/config"
	"5/work/Newyear/user-srv/handler/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Register struct{}
type Login struct{}
type UpdatePassword struct {
}
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
	return config.DB.Where("username=?", username).Update("password", newPassword).Error
}

// 用户查询
func (l *ListDao) ListUser(db *gorm.DB) ([]model.User, error) {
	var users []model.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// 商品添加
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

// 查询所有商品
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

// 查询订单
func (o *OrderDao) GetOrder(payStatus int) (*model.OrderMain, error) {
	var order model.OrderMain
	if err := config.DB.Where("pay_status=?", payStatus).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// 标记订单已支付（根据订单号）
func (o *OrderDao) MarkPaidByOrderNo(orderNo string, tradeNo string) error {
	return config.DB.Model(&model.OrderMain{}).
		Where("order_no=?", orderNo).
		Updates(map[string]interface{}{
			"pay_status": 1,
			"trade_no":   tradeNo,
		}).Error
}

// 单商品立即下单
func (o *OrderDao) CreateOrderWithTx(userID uint, productID uint, buyQty int) (*model.OrderMain, error) {
	var created *model.OrderMain
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var product model.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, productID).Error; err != nil {
			return err
		}
		if product.ProductNum < buyQty {
			return fmt.Errorf("库存不足")
		}
		if err := tx.Model(&product).Update("product_num", product.ProductNum-buyQty).Error; err != nil {
			return err
		}
		order := model.OrderMain{
			OrderNo:     fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
			UserID:      userID,
			TotalAmount: float64(product.ProductPrice) * float64(buyQty),
			OrderStatus: 1,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		item := model.OrderItem{
			OrderID:      order.ID,
			ProductID:    product.ID,
			ProductName:  product.ProductName,
			ProductPrice: product.ProductPrice,
			Quantity:     buyQty,
			Amount:       float64(product.ProductPrice) * float64(buyQty),
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		created = &order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// 从购物车结算
func (o *OrderDao) CreateOrderFromCartTx(userID uint) (*model.OrderMain, error) {
	var created *model.OrderMain
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var carts []model.ShoppingCart
		if err := tx.Where("user_id=?", userID).Find(&carts).Error; err != nil {
			return err
		}
		if len(carts) == 0 {
			return fmt.Errorf("购物车为空")
		}
		total := 0.0
		products := make(map[uint]model.Product)
		for _, c := range carts {
			var p model.Product
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&p, c.ProductID).Error; err != nil {
				return err
			}
			if p.ProductNum < c.Quantity {
				return fmt.Errorf("库存不足: productID=%d", c.ProductID)
			}
			products[c.ProductID] = p
			total += float64(p.ProductPrice) * float64(c.Quantity)
		}
		for _, c := range carts {
			p := products[c.ProductID]
			if err := tx.Model(&p).Update("product_num", p.ProductNum-c.Quantity).Error; err != nil {
				return err
			}
		}
		order := model.OrderMain{
			OrderNo:     fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
			UserID:      userID,
			TotalAmount: total,
			OrderStatus: 1,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for _, c := range carts {
			p := products[c.ProductID]
			item := model.OrderItem{
				OrderID:      order.ID,
				ProductID:    p.ID,
				ProductName:  p.ProductName,
				ProductPrice: p.ProductPrice,
				Quantity:     c.Quantity,
				Amount:       float64(p.ProductPrice) * float64(c.Quantity),
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("user_id=?", userID).Delete(&model.ShoppingCart{}).Error; err != nil {
			return err
		}
		created = &order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// 购物车系列
func (d *CartDao) Add(userID uint, productID uint, qty int) error {
	var item model.ShoppingCart
	err := config.DB.Where("user_id=? and product_id=?", userID, productID).First(&item).Error
	if err == nil {
		item.Quantity = qty
		return config.DB.Save(&item).Error
	}
	newItem := model.ShoppingCart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  qty,
	}
	return config.DB.Create(&newItem).Error
}

func (d *CartDao) List(userID uint) ([]model.ShoppingCart, error) {
	var items []model.ShoppingCart
	if err := config.DB.Where("user_id=?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *CartDao) Remove(userID uint, productID uint) error {
	return config.DB.Where("user_id=? and product_id=?", userID, productID).Delete(&model.ShoppingCart{}).Error
}

func (d *CartDao) Clear(userID uint) error {
	return config.DB.Where("user_id=?", userID).Delete(&model.ShoppingCart{}).Error
}

// 商家入驻审核资质系列
func (a *AuthorAuthDao) Create(auth *model.AuthorAuth) error {
	return config.DB.Create(&auth).Error
}

func (a *AuthorAuthDao) Get(userID uint) (*model.AuthorAuth, error) {
	var auth model.AuthorAuth
	if err := config.DB.Where("user_id=?", userID).First(&auth).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

func (a *AuthorAuthDao) Update(authID uint, status int, rejectReason, auditor string) error {
	return config.DB.Model(&model.AuthorAuth{}).Where("id=?", authID).
		Updates(map[string]interface{}{
			"status":        status,
			"reject_reason": rejectReason,
			"auditor":       auditor,
		}).Error
}
