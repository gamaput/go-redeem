package model

import "gorm.io/gorm"

// User merupakan model untuk data pengguna (user)
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

// TableName mengembalikan nama tabel untuk model User
func (User) TableName() string {
	return "users"
}
