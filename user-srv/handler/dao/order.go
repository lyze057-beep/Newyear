package dao

import (
	"5/work/Newyear/user-srv/basic/config"
	"5/work/Newyear/user-srv/handler/model"
	"errors"

	"gorm.io/gorm"
)

type Register struct{}
type Login struct{}

func (r *Register) Create(user *model.User) error {
	return config.DB.Create(&user).Error
}
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
