package database

import "gorm.io/gorm"

type Models struct {
	Users      UserModel
	Carts      CartModel
	CartItems  CartItemModel
	Orders     OrderModel
	OrderItems OrderItemModel
	Products   ProductModel
	Categories CategoryModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:      UserModel{DB: db},
		Carts:      CartModel{DB: db},
		CartItems:  CartItemModel{DB: db},
		Orders:     OrderModel{DB: db},
		OrderItems: OrderItemModel{DB: db},
		Products:   ProductModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
