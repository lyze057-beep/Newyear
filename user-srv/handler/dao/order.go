package dao

import (
	"5/work/Newyear/user-srv/basic/config"
	"5/work/Newyear/user-srv/handler/model"
	"errors"

	"gorm.io/gorm"
)

type Register struct{}
type Login struct{}

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
type UpdatePassword struct {
}

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

type ListDao struct{}

func (l *ListDao) ListUser(db *gorm.DB) ([]model.User, error) {
	var users []model.User
	config.DB.Find(&users)
	return users, nil
}

type Product struct{}

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
