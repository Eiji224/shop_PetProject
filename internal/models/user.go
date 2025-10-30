package models

type User struct {
	ID       uint   `gorm:"primaryKey;AUTO_INCREMENT"`
	Username string `gorm:"size:100;not null"`
	Email    string `gorm:"uniqueIndex;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	Type     string `gorm:"size:100;not null"`

	Cart *Cart `gorm:"constraint:OnDelete:CASCADE;"`
}
