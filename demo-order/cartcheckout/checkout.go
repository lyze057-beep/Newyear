package cartcheckout

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID    uint
	Name  string
	Stock int
	Price float64
}

type Order struct {
	ID          uint
	OrderNo     string
	UserID      uint
	TotalAmount float64
}

type OrderItem struct {
	ID        uint
	OrderID   uint
	ProductID uint
	Quantity  int
	Price     float64
}

type ShoppingCart struct {
	ID        uint
	UserID    uint
	ProductID uint
	Quantity  int
}

func CheckoutFromCartTx(db *gorm.DB, userID uint) (*Order, error) {
	var created *Order
	err := db.Transaction(func(tx *gorm.DB) error {
		var carts []ShoppingCart
		if err := tx.Where("user_id=?", userID).Find(&carts).Error; err != nil {
			return err
		}
		if len(carts) == 0 {
			return fmt.Errorf("购物车为空")
		}
		total := 0.0
		products := make(map[uint]Product)
		for _, c := range carts {
			var p Product
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&p, c.ProductID).Error; err != nil {
				return err
			}
			if p.Stock < c.Quantity {
				return fmt.Errorf("库存不足: productID=%d", c.ProductID)
			}
			products[c.ProductID] = p
			total += p.Price * float64(c.Quantity)
		}
		for _, c := range carts {
			p := products[c.ProductID]
			if err := tx.Model(&p).Update("stock", p.Stock-c.Quantity).Error; err != nil {
				return err
			}
		}
		order := Order{
			OrderNo:     fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
			UserID:      userID,
			TotalAmount: total,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for _, c := range carts {
			p := products[c.ProductID]
			item := OrderItem{
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  c.Quantity,
				Price:     p.Price,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("user_id=?", userID).Delete(&ShoppingCart{}).Error; err != nil {
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
