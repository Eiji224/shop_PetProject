package database

import (
	"context"
	"time"

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
}

func (u *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := u.DB.WithContext(ctx).Create(user)
	return result.Error
}

func (u *UserModel) FindByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	result := u.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	return &user, result.Error
}
