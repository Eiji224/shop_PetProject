package models

type Category struct {
	ID   uint   `gorm:"primaryKey;AUTO_INCREMENT"`
	Name string `gorm:"size:100;not null"`
}
