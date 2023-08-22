package model

import "gorm.io/gorm"

// Prize adalah model untuk menyimpan informasi hadiah
type Prize struct {
	gorm.Model
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// TableName mengembalikan nama tabel untuk model Prize
func (Prize) TableName() string {
	return "prizes"
}
