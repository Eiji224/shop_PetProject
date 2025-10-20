package database

import (
	"context"

	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

type User struct {
	ID       uint   `gorm:"primaryKey;AUTO_INCREMENT"`
	Username string `gorm:"size:100;not null"`
	Email    string `gorm:"uniqueIndex;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	Type     string `gorm:"size:100;not null"`

	Cart *Cart `gorm:"constraint:OnDelete:CASCADE;"`
}

func (u *UserModel) Insert(user *User, ctx context.Context) error {
	result := u.DB.WithContext(ctx).Create(user)
	return result.Error
}

func (u *UserModel) FindById(id uint, ctx context.Context) (*User, error) {
	var user User
	result := u.DB.WithContext(ctx).First(&user, id)
	return &user, result.Error
}

func (u *UserModel) FindByUsername(username string, ctx context.Context) (*User, error) {
	var user User
	result := u.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	return &user, result.Error
}
