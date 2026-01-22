package models

import "time"

type AuthUser struct {
	Id       string `gorm:"type:varchar(191);primaryKey"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type RefreshToken struct {
	Id        string    `gorm:"type:varchar(191);primaryKey"`
	UserID    string    `gorm:"type:varchar(191);not null"`
	ExpiresAt time.Time `json:"expires_at"`
	User      AuthUser  `gorm:"foreignKey:UserID;references:Id"`
}
