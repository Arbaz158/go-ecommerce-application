package models

type UserProfile struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Phone string `gorm:"size:15;unique;not null"`
	Email string `gorm:"size:100;unique;not null"`
}

type Address struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null; foreignKey:UserProfile(ID)"`
	Street     string `gorm:"size:200;not null"`
	City       string `gorm:"size:100;not null"`
	State      string `gorm:"size:100;not null"`
	PostalCode string `gorm:"size:20;not null"`
}
