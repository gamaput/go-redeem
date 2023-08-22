package model

import "gorm.io/gorm"

type RedeemCode struct {
	gorm.Model
	Code       string `json:"code"`
	IsRedeemed bool   `json:"is_redeemed"`
	PrizeID    uint   `json:"prize_id"` // Kunci asing ke model Prize
	Name       string `json:"name"`
	NoKTP      string `json:"no_ktp"`
	City       string `json:"city"`
	Address    string `json:"address"`
	PhoneNo    string `json:"phone_no"`
}

func (RedeemCode) TableName() string {
	return "redeem_codes"
}
